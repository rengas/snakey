package game

type Fruit struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewFruit(x, y int) Fruit {
	return Fruit{X: x, Y: y}
}

func (f Fruit) IsEaten(x, y int) bool {
	return f.X == x && f.Y == y
}
