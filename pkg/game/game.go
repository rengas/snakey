package game

import (
	"errors"
	"fmt"
)

var (
	ErrOutOfArena      = errors.New("out of arena")
	ErrInvalidMovement = errors.New("invalid movement")
	ErrIsNotEaten      = errors.New("not eaten")
)

type Game struct {
	Arena Arena `json:"arena"`
	Fruit Fruit `json:"fruit"`
	Snake Snake `json:"snake"`
}

// New ...
func New(arena Arena, snake Snake, fruit Fruit) Game {
	return Game{
		Arena: arena,
		Fruit: fruit,
		Snake: snake,
	}
}

func (g Game) Scored(posX, posY int) error {
	if !g.Arena.IsInside(posX, posY) {
		fmt.Println(posX, posY)
		return ErrOutOfArena
	}

	if !g.Snake.CanMove(posX, posY) {
		return ErrInvalidMovement
	}

	if !g.Fruit.IsEaten(posX, posY) {
		return ErrIsNotEaten
	}

	return nil
}
