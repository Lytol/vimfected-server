package server

import (
	"encoding/json"

	"github.com/Lytol/vimfected-server/game"
)

type CommandType string

const (
	Register CommandType = "register"
	Snapshot CommandType = "snapshot"
)

type Command struct {
	Type CommandType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type RegisterData struct {
	Id string `json:"id"`
}

type SnapshotData struct {
	Players []*game.Player `json:"players"`
	Map     *game.Map      `json:"map"`
}

func SnapshotCommand(g *game.Game) (Command, error) {
	players := make([]*game.Player, len(g.Players))
	i := 0
	for _, player := range g.Players {
		players[i] = player
		i++
	}

	snapshotData := SnapshotData{
		Players: players,
		Map:     g.Map,
	}

	data, err := json.Marshal(snapshotData)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: Snapshot,
		Data: data,
	}, nil
}
