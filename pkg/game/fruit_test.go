package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEaten(t *testing.T) {
	type test struct {
		name   string
		inputX int
		inputY int
		want   bool
	}

	fruit := Fruit{
		X: 5,
		Y: 5,
	}

	tests := []test{
		{name: "given X and Y are same as fruit location, then return true", inputX: 5, inputY: 5, want: true},
		{name: "given Y is same as fruit y but different X, then should return false", inputX: 0, inputY: 5, want: false},
		{name: "given X is same as fruit X but different Y, then should return false", inputX: 5, inputY: 2, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, fruit.IsEaten(tc.inputX, tc.inputY))
		})
	}
}
