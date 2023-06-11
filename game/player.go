package game

import "time"

const (
	PlayerSpeed = 100
)

type Player struct {
	Id string `json:"id"`
	X  int64  `json:"x"`
	Y  int64  `json:"y"`

	lastMovedAt time.Time
}

func (p *Player) Move(x int64, y int64) {
	p.X = x
	p.Y = y
	p.lastMovedAt = time.Now()
}

func (p *Player) IsReadyToMove() bool {
	return time.Since(p.lastMovedAt) > time.Millisecond*PlayerSpeed
}
