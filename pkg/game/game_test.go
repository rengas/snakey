package game

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestScored(t *testing.T) {
	type test struct {
		name    string
		game    Game
		posX    int
		posY    int
		wantErr error
	}

	width := 20
	height := 20

	tests := []test{
		{
			name: "given the snake is inside arena and try to move outside of area, then an error should be returned",
			game: Game{
				Arena: Arena{
					Width:  width,
					Height: height,
				},
				Fruit: NewFruit(rand.Intn(width-1)+1, rand.Intn(height-1)+1),
				Snake: NewSnake(5, 5),
			},
			posX:    -1,
			posY:    -1,
			wantErr: ErrOutOfArena,
		},
		{
			name: "given the snake is inside arena, and try to 180 degree turn, then error should be thrown",
			game: Game{
				Arena: Arena{
					Width:  width,
					Height: height,
				},
				Fruit: NewFruit(rand.Intn(width-1)+1, rand.Intn(height-1)+1),
				Snake: NewSnake(1, 1),
			},
			posX:    1,
			posY:    1,
			wantErr: ErrInvalidMovement,
		},
		{
			name: "given the snake and fruit location is inside arena,but snake didn't eat the fruit ,then an error should be throw ",
			game: Game{
				Arena: Arena{
					Width:  width,
					Height: height,
				},
				Fruit: NewFruit(1, 2),
				Snake: NewSnake(1, 4),
			},
			posX:    1,
			posY:    1,
			wantErr: ErrIsNotEaten,
		},
		{
			name: "given the snake and fruit location is inside arena,but snake didn't eat the fruit ,then an error should be throw ",
			game: Game{
				Arena: Arena{
					Width:  width,
					Height: height,
				},
				Fruit: NewFruit(1, 4),
				Snake: NewSnake(1, 3),
			},
			posX:    1,
			posY:    4,
			wantErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.game.Scored(tc.posX, tc.posY)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
