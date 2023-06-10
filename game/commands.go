package game

import (
	"encoding/json"

	"github.com/Lytol/vimfected-server/commands"
)

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

func (g *Game) AddPlayerCommand(player *Player) (commands.Command, error) {
	data, err := json.Marshal(player)
	if err != nil {
		return commands.Command{}, err
	}

	return commands.Command{
		Type: commands.AddPlayer,
		Id:   player.Id,
		Data: data,
	}, nil
}
