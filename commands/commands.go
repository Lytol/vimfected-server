package commands

import (
	"encoding/json"
)

type Type string

const (
	Register      Type = "register"
	Snapshot      Type = "snapshot"
	SpawnPlayer   Type = "spawn_player"
	MoveDirection Type = "move_direction"
	MovePosition  Type = "move_position"
)

type Command struct {
	Type Type            `json:"type"`
	Id   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

type RegisterData struct{}

func SpawnPlayerCommand(id string) (Command, error) {
	return Command{
		Type: SpawnPlayer,
		Id:   id,
	}, nil
}

type MoveDirectionData struct {
	Direction string `json:"direction"`
}

type MovePositionData struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}
