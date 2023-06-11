package game

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/Lytol/vimfected-server/commands"
)

const (
	TickRate            = 60
	CommandOutputBuffer = 100
)

type Game struct {
	Players map[string]*Player
	Map     *Map

	Incoming *commands.Queue
	Outgoing chan commands.Command

	quit chan bool
}

func New() (*Game, error) {
	var err error

	g := &Game{
		Players: make(map[string]*Player),
		quit:    make(chan bool),

		Incoming: commands.NewQueue(),
		Outgoing: make(chan commands.Command, CommandOutputBuffer),
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

type Notify func(cmd commands.Command) error

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

func (g *Game) Queue(cmd commands.Command) {
	g.Incoming.Shift(cmd)
}

func (g *Game) Update(delta float64) {
	// Process all pending commands
	for {
		cmd, ok := g.Incoming.Unshift()
		if !ok {
			return
		}
		log.Printf("Incoming Command %s | Id: %s | %s\n", cmd.Type, cmd.Id, cmd.Data)
		g.handleCommand(cmd)
	}
}

func (g *Game) handleCommand(cmd commands.Command) {
	switch cmd.Type {
	case commands.SpawnPlayer:
		g.Outgoing <- cmd
	default:
		log.Printf("Unknown command: %s\n", cmd.Type)
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

	cmd, err := g.AddPlayerCommand(player)
	if err != nil {
		return nil, err
	}

	g.Outgoing <- cmd

	return player, nil
}

func (g *Game) AddPlayer(p *Player) error {
	g.Players[p.Id] = p
	log.Printf("Added player: %s\n", p.Id)
	return nil
}

func (g *Game) RemovePlayer(p *Player) error {
	if _, ok := g.Players[p.Id]; ok {
		delete(g.Players, p.Id)
	} else {
		return fmt.Errorf("cannot remove player, does not exist: %s", p.Id)
	}
	log.Printf("Removed player: %s\n", p.Id)
	return nil
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
