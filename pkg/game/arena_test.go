package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsInside(t *testing.T) {
	type test struct {
		name string

		inputX int
		inputY int
		want   bool
	}

	arena := Arena{
		Width:  5,
		Height: 5,
	}

	tests := []test{
		{name: "given X and Y on arena border, then should return false", inputX: 0, inputY: 0, want: false},
		{name: "given Y on the line, then should return false", inputX: 0, inputY: 5, want: false},
		{name: "given X on the line, then should return false", inputX: 5, inputY: 0, want: false},
		{name: "given X and Y outside arena, then should return false", inputX: 6, inputY: 6, want: false},
		{name: "given X  outside arena but y inside, then should return false", inputX: 6, inputY: 3, want: false},
		{name: "given X  inside arena but y outside, then should return false", inputX: 3, inputY: 6, want: false},
		{name: "given X and Y are outside, then should return false", inputX: -1, inputY: -4, want: false},
		{name: "given X and Y are inside, then should return false", inputX: 1, inputY: 4, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, arena.IsInside(tc.inputX, tc.inputY))
		})
	}

}
