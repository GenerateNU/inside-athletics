package goat

import (
	"context"
	"inside-athletics/internal/utils"
)

type GoatService struct {
	goatDB *GoatDB // give the service a DB connection!
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

func (g *GoatService) GetGoat(ctx context.Context, input *GetGoatParams) (*utils.ResponseBody[GoatResponse], error) {
	resp := &utils.ResponseBody[GoatResponse]{}
	id := input.Id
	goat, err := g.goatDB.GetGoat(uint(id))

	if err != nil { // if we have an err getting the goat return the empty response and the error
		return resp, err
	}

	resp.Body = &GoatResponse{
		Id:   int8(id),
		Name: goat.Name,
		Age:  goat.Age,
	}

	return resp, nil
}
