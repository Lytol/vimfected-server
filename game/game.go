package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"
)

const (
	TickRate            = 60
	CommandOutputBuffer = 100
)

type Game struct {
	Players map[string]*Player
	Map     *Map

	Incoming *Queue
	Outgoing chan Command

	quit chan bool
}

func New() (*Game, error) {
	var err error

	g := &Game{
		Players: make(map[string]*Player),
		quit:    make(chan bool),

		Incoming: NewQueue(),
		Outgoing: make(chan Command, CommandOutputBuffer),
	}

	g.Map, err = NewMap(DefaultMapWidth, DefaultMapHeight)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) Run() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	tickInterval := time.Second / TickRate
	timeStart := time.Now().UnixNano()

	ticker := time.NewTicker(tickInterval)

	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			// DT in seconds
			delta := float64(now-timeStart) / 1000000000
			timeStart = now
			g.Update(delta)
		case <-g.quit:
			ticker.Stop()
			return
		}
	}
}

type Notify func(cmd Command) error

func (g *Game) Notifier(notify Notify) {
	for {
		select {
		case cmd := <-g.Outgoing:
			log.Printf("Outgoing Command %s | Id: %s | %s\n", cmd.Type, cmd.Id, cmd.Data)
			err := notify(cmd)
			if err != nil {
				log.Printf("Notifier error: %v\n", err)
			}
		case <-g.quit:
			return
		}
	}
}

func (g *Game) Stop() {
	g.quit <- true
}

func (g *Game) Queue(cmd Command) {
	g.Incoming.Shift(cmd)
}

func (g *Game) Update(delta float64) {
	// Process all pending commands
	g.Incoming.Each(func(cmd Command) bool {
		return g.handleCommand(cmd)
	})
}

func (g *Game) handleCommand(cmd Command) bool {
	switch cmd.Type {
	case ClearPlayerInput:
		g.Incoming.ClearUntil(cmd)
		return true
	case MovePlayerInput:
		var movePlayerInputData MovePlayerInputData

		// TODO: players should only be able to move themselves
		player, ok := g.Players[cmd.Id]
		if !ok {
			log.Printf("Player not found: %s\n", cmd.Id)
			return true
		}

		if !player.IsReady() {
			return false
		}

		err := json.Unmarshal(cmd.Data, &movePlayerInputData)
		if err != nil {
			log.Printf("Error unmarshalling move player input data: %v\n", err)
			return true
		}

		g.MovePlayer(player, movePlayerInputData.Direction)

		return true
	default:
		log.Printf("Unknown command: %s\n", cmd.Type)
		return true
	}
}

func (g *Game) SpawnPlayer(id string) (*Player, error) {
	x, y := g.findSpawn()

	player := &Player{
		Id: id,
		X:  x,
		Y:  y,
	}

	err := g.AddPlayer(player)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (g *Game) AddPlayer(p *Player) error {
	g.Players[p.Id] = p
	log.Printf("Added player: %s\n", p.Id)

	cmd, err := AddPlayerCommand(p)
	if err != nil {
		return err
	}

	g.Outgoing <- cmd

	return nil
}

func (g *Game) RemovePlayer(p *Player) error {
	if _, ok := g.Players[p.Id]; ok {
		delete(g.Players, p.Id)
	} else {
		return fmt.Errorf("cannot remove player, does not exist: %s", p.Id)
	}
	log.Printf("Removed player: %s\n", p.Id)

	cmd, err := RemovePlayerCommand(p)
	if err != nil {
		return err
	}

	g.Outgoing <- cmd

	return nil
}

func (g *Game) MovePlayer(player *Player, dir Direction) bool {
	newX := player.X
	newY := player.Y

	switch dir {
	case DirectionUp:
		newY--
	case DirectionDown:
		newY++
	case DirectionLeft:
		newX--
	case DirectionRight:
		newX++
	}

	if newX < 0 || newX >= g.Map.Width || newY < 0 || newY >= g.Map.Height {
		return false
	}

	if g.playerAt(newX, newY) != nil {
		return false
	}

	player.Move(newX, newY)

	cmd, err := MovePlayerCommand(player)
	if err != nil {
		log.Printf("Error creating move player command: %v\n", err)
		return false
	}

	g.Outgoing <- cmd

	return true
}

func (g *Game) findSpawn() (int64, int64) {
	for {
		x := rand.Int63n(g.Map.Width)
		y := rand.Int63n(g.Map.Height)

		if g.playerAt(x, y) == nil {
			return x, y
		}
	}
}

func (g *Game) playerAt(x int64, y int64) *Player {
	for _, player := range g.Players {
		if player.X == x && player.Y == y {
			return player
		}
	}
	return nil
}
