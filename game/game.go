package game

import (
	"fmt"
	"log"
)

type Game struct {
	Players map[string]*Player
}

func New() *Game {
	return &Game{
		Players: make(map[string]*Player),
	}
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
