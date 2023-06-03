package server

type CommandType int64

const (
	Register CommandType = 1
)

type Command struct {
	Type CommandType
	Data map[string]string
}
