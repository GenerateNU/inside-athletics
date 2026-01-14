package goat

import (
	"context"
	"inside-athletics/internal/utils"
)

type GoatService struct {
}

func (g *GoatService) Ping(ctx context.Context, input *utils.EmptyInput) (*utils.ResponseBody[GoatResponse], error) {
	resp := &GoatResponse{
		Id:   1,
		Name: "Suli",
		Age:  67,
	}
	return &utils.ResponseBody[GoatResponse]{
		Body: resp,
	}, nil
}
