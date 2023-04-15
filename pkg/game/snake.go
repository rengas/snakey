package game

type Snake struct {
	X int `json:"currentX"`
	Y int `json:"currentY"`
}

func NewSnake(x int, y int) Snake {
	return Snake{
		X: x,
		Y: y,
	}
}

func (g Snake) CanMove(posX, posY int) bool {
	return g.X != posX || g.Y != posY
}
