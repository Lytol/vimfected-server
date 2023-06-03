package game

import (
	"fmt"
	"log"
)

type Game struct {
	Players map[string]*Player
	Map     *Map
}

func New() (*Game, error) {
	var err error

	g := &Game{
		Players: make(map[string]*Player),
	}

	g.Map, err = NewMap(DefaultMapWidth, DefaultMapHeight)
	if err != nil {
		return nil, err
	}

	return g, nil
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

func (g *Game) Snapshot() *Snapshot {
	snapshot := &Snapshot{
		Players: make([]*Player, len(g.Players)),
		Map:     g.Map,
	}

	i := 0
	for _, player := range g.Players {
		snapshot.Players[i] = player
		i += 1
	}

	return snapshot
}

type Snapshot struct {
	Players []*Player
	Map     *Map
}
