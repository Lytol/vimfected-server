package game

import "time"

const (
	PlayerSpeed = 250 // ms per action
)

type Player struct {
	Id string `json:"id"`
	X  int64  `json:"x"`
	Y  int64  `json:"y"`

	lastActionAt time.Time
}

func (p *Player) Move(x int64, y int64) {
	p.X = x
	p.Y = y
	p.lastActionAt = time.Now()
}

func (p *Player) IsReady() bool {
	return time.Since(p.lastActionAt) > time.Millisecond*PlayerSpeed
}
