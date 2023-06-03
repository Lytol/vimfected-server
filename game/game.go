package game

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/Lytol/vimfected-server/commands"
)

const (
	TickRate = 60
)

type Game struct {
	Commands *list.List
	Players  map[string]*Player
	Map      *Map

	quit chan bool
}

func New() (*Game, error) {
	var err error

	g := &Game{
		Commands: list.New(),
		Players:  make(map[string]*Player),
		quit:     make(chan bool),
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

func (g *Game) Stop() {
	g.quit <- true
}

func (g *Game) Update(delta float64) {
	// Process all pending commands
	for {
		cmd, ok := g.Dequeue()
		if !ok {
			return
		}
		log.Printf("Command %s | Id: %s | %s\n", cmd.Type, cmd.Id, cmd.Data)
	}
}

func (g *Game) Queue(cmd commands.Command) {
	g.Commands.PushBack(cmd)
}

func (g *Game) Dequeue() (commands.Command, bool) {
	if g.Commands.Len() == 0 {
		return commands.Command{}, false
	}
	return g.Commands.Remove(g.Commands.Front()).(commands.Command), true
}

func (g *Game) CreatePlayer(id string) (*Player, error) {
	p := &Player{
		Id: id,
	}
	g.Players[id] = p
	log.Printf("Created player: %s\n", p.Id)
	return p, nil
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

type snapshotData struct {
	Players []*Player `json:"players"`
	Map     *Map      `json:"map"`
}

func (g *Game) SnapshotCommand() (commands.Command, error) {
	players := make([]*Player, len(g.Players))
	i := 0
	for _, player := range g.Players {
		players[i] = player
		i++
	}

	snapshotData := snapshotData{
		Players: players,
		Map:     g.Map,
	}

	data, err := json.Marshal(snapshotData)
	if err != nil {
		return commands.Command{}, err
	}

	return commands.Command{
		Type: commands.Snapshot,
		Data: data,
	}, nil
}
