package game

const (
	TileWidth        = 16
	TileHeight       = 16
	DefaultMapWidth  = 500
	DefaultMapHeight = 500
)

const (
	GrassDefault = 317
)

type Map struct {
	Width  int64     `json:"width"`
	Height int64     `json:"height"`
	Data   [][]int64 `json:"data"`
}

func NewMap(width int64, height int64) (*Map, error) {
	m := &Map{
		Width:  width,
		Height: height,
		Data:   make([][]int64, width),
	}

	for i := range m.Data {
		m.Data[i] = make([]int64, height)
		for j := range m.Data[i] {
			m.Data[i][j] = GrassDefault
		}
	}

	return m, nil
}
