package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	skerror "snakey/pkg/errors"
	"snakey/pkg/game"
	"snakey/pkg/httputils"
)

const (
	ErrRestaurantInvalidWidthOrHeight skerror.ValidationError = "invalid width or height"
	ErrTicksDontLeadToFruit           skerror.ValidationError = "don't lead to fruit"
	ErrOutOfArena                     skerror.ValidationError = "out of arena"
	ErrInvalidMovement                skerror.ValidationError = "invalid movement"
)

type SnakeyAPI struct{}

func NewSnakeyAPI() SnakeyAPI {
	return SnakeyAPI{}
}

type Game struct {
	GameID string `json:"gameId"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Score  int    `json:"score"`
	Fruit  fruit  `json:"fruit"`
	Snake  snake  `json:"snake"`
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

type NewResponse struct {
	Game
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (a *SnakeyAPI) New(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	width, height, err := httputils.GetWidthAndHeight(req)
	if err != nil {
		if err != nil {
			httputils.UnProcessableEntity(ctx, w, ErrRestaurantInvalidWidthOrHeight)
			return
		}
	}

	httputils.OK(ctx, w, NewResponse{
		Game: Game{
			GameID: uuid.NewString(),
			Width:  width,
			Height: height,
			Score:  0,
			Fruit: fruit{
				X: rand.Intn(width-1) + 1,
				Y: rand.Intn(height-1) + 1,
			},
			Snake: snake{
				X: 1,
				Y: 1,
			},
		},
	})

}

type Tick struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ValidateRequest struct {
	Game  Game   `json:"game"`
	Ticks []Tick `json:"ticks"`
}

type ValidateReponse struct {
	Game Game `json:"game"`
}

func (a *SnakeyAPI) Validate(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var ur ValidateRequest
	err := httputils.ReadJson(req, &ur)
	if err != nil {
		httputils.BadRequest(ctx, w, err)
		return
	}

	currentScore := ur.Game.Score
	ar := game.NewArena(ur.Game.Width, ur.Game.Height)
	sn := game.NewSnake(ur.Game.Snake.X, ur.Game.Snake.Y)
	fr := game.NewFruit(ur.Game.Fruit.X, ur.Game.Fruit.Y)
	state := game.New(ar, sn, fr)
	for _, v := range ur.Ticks {
		newPosX := state.Snake.X + v.X
		newPosY := state.Snake.Y + v.Y
		log.Printf("snake x :%d, snake y:%d", newPosX, newPosY)
		err := state.Scored(newPosX, newPosY)

		if err != nil {
			if errors.Is(err, game.ErrOutOfArena) {
				httputils.TeaPot(ctx, w, ErrOutOfArena)
				return
			} else if errors.Is(err, game.ErrInvalidMovement) {
				httputils.TeaPot(ctx, w, ErrInvalidMovement)
				return
			}
		}

		if errors.Is(err, game.ErrIsNotEaten) {
			state.Snake.X = newPosX
			state.Snake.Y = newPosY
			continue
		}
		currentScore += 1
	}

	if currentScore == ur.Game.Score {
		httputils.BadRequest(ctx, w, ErrTicksDontLeadToFruit)
		return
	}
	ur.Game.Score = currentScore
	ur.Game.Snake.X = state.Snake.X
	ur.Game.Snake.Y = state.Snake.Y
	httputils.OK(ctx, w, ValidateReponse{
		Game: ur.Game,
	})
}

func (a *SnakeyAPI) Health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "I'm ok")
}
