package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/Lytol/vimfected-server/game"
)

type Server struct {
	Game *game.Game

	serveMux http.ServeMux
}

func NewServer() (*Server, error) {
	var err error

	s := &Server{}

	s.Game, err = game.New()
	if err != nil {
		return nil, err
	}

	s.serveMux.HandleFunc("/", s.handle)

	return s, nil
}

func (s *Server) Run() error {
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
		cmd    Command
		err    error
		player *game.Player
	)

	for {
		err = wsjson.Read(ctx, ws, &cmd)
		if err != nil {
			return err
		}

		switch cmd.Type {
		case Register:
			var data RegisterData

			err = json.Unmarshal(cmd.Data, &data)
			if err != nil {
				return ws.Close(websocket.StatusProtocolError, "Invalid data for register")
			}

			player, err = s.Game.CreatePlayer(data.Id)
			if err != nil {
				return ws.Close(websocket.StatusProtocolError, "Unable to register player")
			}
			defer s.Game.RemovePlayer(player)

			snapshot, err := SnapshotCommand(s.Game)
			if err != nil {
				return ws.Close(websocket.StatusProtocolError, "Unable to create snapshot")
			}

			err = wsjson.Write(ctx, ws, snapshot)
			if err != nil {
				return ws.Close(websocket.StatusProtocolError, "Unable to write snapshot")
			}
		default:
			log.Printf("Command | type: %s | data: %+v\n", cmd.Type, cmd.Data)
		}
	}
}
