package game

type Arena struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewArena(width, height int) Arena {
	return Arena{
		Width:  width,
		Height: height,
	}
}

func (a Arena) IsInside(x, y int) bool {
	return x > 0 && x < a.Width && y > 0 && y < a.Height
}
