package httputils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Error error `json:"error"`
}

func OK(ctx context.Context, w http.ResponseWriter, v interface{}) {
	WriteJSON(ctx, w, v, http.StatusOK)
}

func UnProcessableEntity(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusUnprocessableEntity)
}

func TeaPot(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusTeapot)
}

func BadRequest(ctx context.Context, w http.ResponseWriter, err error) {
	WriteJSON(ctx, w, ErrorResponse{Error: err}, http.StatusBadRequest)
}

func ReadJson(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, v interface{}, statusCode int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}

func GetWidthAndHeight(req *http.Request) (int, int, error) {
	width := 10
	height := 1

	if req.URL.Query().Get("w") != "" {
		p, err := strconv.Atoi(req.URL.Query().Get("w"))
		if err != nil {
			return 0, 0, err
		}
		width = p
	}
	if req.URL.Query().Get("h") != "" {
		l, err := strconv.Atoi(req.URL.Query().Get("h"))
		if err != nil {
			return width, height, err
		}
		height = l
	}

	if width <= 0 && height <= 0 {
		return 0, 0, errors.New("invalid width or height")
	}
	return width, height, nil
}
