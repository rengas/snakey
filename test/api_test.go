//go:build e2e

package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Game struct {
	GameID string `json:"gameId"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Score  int    `json:"score"`
	Fruit  fruit  `json:"fruit"`
	Snake  snake  `json:"snake"`
}
type NewResponse struct {
	Game Game `json:"game"`
}

type fruit struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type snake struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	VelX int `json:"velX"`
	VelY int `json:"velY"`
}

func TestNewApi(t *testing.T) {

	tests := []struct {
		Name       string         `json:"name"`
		Width      string         `json:"width"`
		Height     string         `json:"height"`
		Response   NewResponse    `json:"response"`
		WantErr    *ErrorResponse `json:"err"`
		StatusCode interface{}    `json:"HTTPStatusCode"`
	}{
		{
			Name:       "Give negative Width and height, then return invalid width or height error",
			Width:      "-10",
			Height:     "-1",
			WantErr:    &ErrorResponse{Error: "invalid width or height"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:       "Give invalid Width and height, then return invalid width or height error",
			Width:      "w",
			Height:     "h",
			WantErr:    &ErrorResponse{Error: "invalid width or height"},
			StatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:       "Give valid Width and height, then game object should be returned",
			Width:      "10",
			Height:     "10",
			StatusCode: http.StatusOK,
		},
	}
	for _, ts := range tests {
		t.Run(ts.Name, func(t *testing.T) {

			t.Parallel()
			req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/new?w=%s&h=%s", httpPort, ts.Width, ts.Height), nil)
			require.NoError(t, err)
			client := &http.Client{
				Transport: &http.Transport{},
			}
			resp, err := client.Do(req)

			require.NoError(t, err)
			require.Equal(t, ts.StatusCode, resp.StatusCode)

			if ts.WantErr != nil {
				var e ErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&e)
				require.NoError(t, err)
				require.Equal(t, ts.WantErr.Error, e.Error)
			}
			if ts.WantErr == nil {
				var ls NewResponse
				err = json.NewDecoder(resp.Body).Decode(&ls)
				require.NoError(t, err)
				require.NotEmpty(t, ls)
			}
		})
	}
}

type Tick struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ValidateRequest struct {
	Game  Game   `json:"game"`
	Ticks []Tick `json:"ticks"`
}

type ValidateResponse struct {
	Game Game `json:"game"`
}

func TestValidateApi(t *testing.T) {
	tests := []struct {
		Name       string           `json:"name"`
		Request    ValidateRequest  `json:"request"`
		Response   ValidateResponse `json:"response"`
		WantErr    *ErrorResponse   `json:"err"`
		StatusCode interface{}      `json:"HTTPStatusCode"`
	}{
		{
			Name: "given a range of velocities that collides with border, then it should return out of arena error",
			Request: ValidateRequest{
				Game: Game{
					Width:  10,
					Height: 10,
					Score:  0,
					Fruit: fruit{
						X: 6,
						Y: 9,
					},
					Snake: snake{
						X: 1,
						Y: 1,
					},
				},
				Ticks: []Tick{
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
				},
			},
			WantErr:    &ErrorResponse{Error: "out of arena"},
			StatusCode: http.StatusTeapot,
		},
		{
			Name: "given a range of velocities that are valid, then score should be increased by one",
			Request: ValidateRequest{
				Game: Game{
					Width:  10,
					Height: 10,
					Score:  0,
					Fruit: fruit{
						X: 6,
						Y: 9,
					},
					Snake: snake{
						X: 1,
						Y: 1,
					},
				},
				Ticks: []Tick{
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 1, Y: 0},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
					{X: 0, Y: 1},
				},
			},
			StatusCode: http.StatusOK,
			Response: ValidateResponse{
				Game: Game{
					Width:  10,
					Height: 10,
					Score:  1,
					Fruit: fruit{
						X: 6,
						Y: 9,
					},
					Snake: snake{
						X: 6,
						Y: 9,
					},
				},
			},
		},
	}

	for _, ts := range tests {
		t.Run(ts.Name, func(t *testing.T) {
			b, err := json.Marshal(ts.Request)
			require.NoError(t, err)
			t.Parallel()
			req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/validate", httpPort), bytes.NewReader(b))
			require.NoError(t, err)
			client := &http.Client{
				Transport: &http.Transport{},
			}
			resp, err := client.Do(req)

			require.NoError(t, err)
			require.Equal(t, ts.StatusCode, resp.StatusCode)

			if ts.WantErr != nil {
				var e ErrorResponse
				err = json.NewDecoder(resp.Body).Decode(&e)
				require.NoError(t, err)
				require.Equal(t, ts.WantErr.Error, e.Error)
			}
			if ts.WantErr == nil {
				var ls ValidateResponse
				err = json.NewDecoder(resp.Body).Decode(&ls)
				require.NoError(t, err)
				require.NotEmpty(t, ls.Game)
				require.Equal(t, 1, ts.Response.Game.Score)
				require.Equal(t, ts.Response.Game.Snake.X, ts.Request.Game.Fruit.X)
				require.Equal(t, ts.Response.Game.Snake.Y, ts.Request.Game.Fruit.Y)
			}
		})
	}

}
