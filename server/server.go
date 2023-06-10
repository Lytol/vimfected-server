package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/Lytol/vimfected-server/commands"
	"github.com/Lytol/vimfected-server/game"
)

type Server struct {
	Game        *game.Game
	Subscribers map[string]*websocket.Conn

	serveMux http.ServeMux
}

func NewServer(g *game.Game) (*Server, error) {
	s := &Server{
		Game:        g,
		Subscribers: make(map[string]*websocket.Conn),
	}

	s.serveMux.HandleFunc("/", s.handle)

	return s, nil
}

func (s *Server) Run() error {
	go s.Game.Run()
	go s.Game.Notifier(func(cmd commands.Command) error {
		return s.Broadcast(cmd)
	})

	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		return err
	}
	log.Printf("listening on %v", l.Addr())

	hs := &http.Server{
		Handler:      s,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- hs.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v\n", err)
	case sig := <-sigs:
		log.Printf("terminating: %v\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Printf("shutting down\n")
	s.Game.Stop()
	return hs.Shutdown(ctx)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols:   []string{"vimfected"},
		OriginPatterns: []string{"localhost:5173", "vimfected.com"},
	})
	if err != nil {
		log.Printf("error accepting websocket: %v\n", err)
		return
	}
	defer ws.Close(websocket.StatusInternalError, "")

	if ws.Subprotocol() != "vimfected" {
		ws.Close(websocket.StatusPolicyViolation, "client must speak the vimfected subprotocol")
		return
	}

	err = s.subscribe(r.Context(), ws)
	if errors.Is(err, context.Canceled) {
		log.Printf("context cancelled\n")
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		log.Printf("other: %v\n", err)
		return
	}
}

func (s *Server) subscribe(ctx context.Context, ws *websocket.Conn) error {
	var (
		cmd        commands.Command
		err        error
		registered bool
	)

	registered = false

	for {
		err = wsjson.Read(ctx, ws, &cmd)
		if err != nil {
			return err
		}

		if !registered && cmd.Type == commands.Register {
			log.Printf("subscribing %s\n", cmd.Id)
			s.Subscribers[cmd.Id] = ws
			defer func() {
				log.Printf("unsubscribing %s\n", cmd.Id)
				delete(s.Subscribers, cmd.Id)
			}()

			err = s.Game.SpawnPlayer(cmd.Id)
			if err != nil {
				return err
			}

			snapshot, err := s.Game.SnapshotCommand()
			if err != nil {
				return err
			}
			s.Send(cmd.Id, snapshot)

			registered = true
		} else if !registered {
			return fmt.Errorf("player has not registered")
		} else {
			s.Game.Queue(cmd)
		}
	}
}

func (s *Server) Send(Id string, cmd commands.Command) error {
	subscriber, ok := s.Subscribers[Id]
	if !ok {
		return fmt.Errorf("subscriber does not exist: %s", Id)
	}
	return wsjson.Write(context.TODO(), subscriber, cmd)
}

func (s *Server) Broadcast(cmd commands.Command) error {
	for _, subscriber := range s.Subscribers {
		err := wsjson.Write(context.TODO(), subscriber, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
