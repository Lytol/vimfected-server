package game

import (
	"encoding/json"
	"time"
)

type Type string

const (
	Register     Type = "register"
	Snapshot     Type = "snapshot"
	AddPlayer    Type = "add_player"
	RemovePlayer Type = "remove_player"
	MovePlayer   Type = "move_player"

	MovePlayerInput  Type = "move_player_input"
	ClearPlayerInput Type = "clear_player_input"
)

type Command struct {
	Type      Type            `json:"type"`
	Id        string          `json:"id"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"-"`
}

type RegisterData struct{}

type SnapshotData struct {
	Players []*Player `json:"players"`
	Map     *Map      `json:"map"`
}

func SnapshotCommand(g *Game) (Command, error) {
	players := make([]*Player, len(g.Players))
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

type AddPlayerData struct {
	*Player
}

func AddPlayerCommand(player *Player) (Command, error) {
	addPlayerData := AddPlayerData{player}

	data, err := json.Marshal(addPlayerData)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: AddPlayer,
		Id:   player.Id,
		Data: data,
	}, nil
}

type RemovePlayerData struct {
	*Player
}

func RemovePlayerCommand(player *Player) (Command, error) {
	removePlayerData := RemovePlayerData{player}

	data, err := json.Marshal(removePlayerData)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: RemovePlayer,
		Id:   player.Id,
		Data: data,
	}, nil
}

type MovePlayerInputData struct {
	Direction Direction `json:"direction"`
}

type MovePlayerData struct {
	*Player
}

func MovePlayerCommand(player *Player) (Command, error) {
	movePlayerData := MovePlayerData{player}

	data, err := json.Marshal(movePlayerData)
	if err != nil {
		return Command{}, err
	}

	return Command{
		Type: MovePlayer,
		Id:   player.Id,
		Data: data,
	}, nil
}
