package sport

import (
	"context"
	"inside-athletics/internal/utils"
)

type SportService struct {
	sportDB *SportDB
}

// connect to DB and create a new sport
func (u *SportService) CreateSport(ctx context.Context, input *struct {
	Body CreateSportRequest
}) (*utils.ResponseBody[SportResponse], error) {
	body != &input.Body
	sport, err := u.sportDB.CreateSport(body)
	respBody := &utils.ResponseBody[SportResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &SportResponse{
		ID:   sport.ID,
		Name: sport.Name,
	}

	return &utils.ResponseBody[SportResponse]{
		Body: response,
	}, nil
}

// connect to DB and get a sport by name
func (u *SportService) GetSportByName(ctx context.Context, input *GetSportByNameParams) (*utils.ResponseBody[SportResponse], error) {
	name := input.Name
	sport, err := u.sportDB.GetSportByName(name)
	respBody := &utils.ResponseBody[SportResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &SportResponse{
		ID:   sport.ID,
		Name: sport.Name,
	}

	return &utils.ResponseBody[SportResponse]{
		Body: response,
	}, err
}

// connect to DB and get all sports
func (u *SportService) GetAllSports(ctx context.Context, input *GetAllSportsParams) (*utils.ResponseBody[GetAllSportsResponse], error) {
	sports, err := u.sportDB.GetAllSports()
	respBody := &utils.ResponseBody[GetAllSportsResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetAllSportsResponse{
		Sports: sports,
	}

	return &utils.ResponseBody[GetAllSportsResponse]{
		Body: response,
	}, nil
}

//connect to DB and get sport by ID
func (u *SportService) GetSportByID(ctx context.Context, input *GetSportByIDParams) (*utils.ResponseBody[SportResponse], error) {
	id := input.ID
	sport, err := u.sportDB.GetSportById(id)
	respBody := &utils.ResponseBody[SportResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &SportResponse{
		ID:   sport.ID,
		Name: sport.Name,
	}

	return &utils.ResponseBody[SportResponse]{
		Body: response,
	}, nil
}

// Connect to DB and update existing sport
func (u *SportService) UpdateSport(ctx context.Context, input *struct {
    Body UpdateSportRequest 
}) (*utils.ResponseBody[SportResponse], error) {
    body := &input.Body
    sport, err := u.sportDB.UpdateSport(body.ID, body.Name)
    respBody := &utils.ResponseBody[SportResponse]{}

    if err != nil {
        return respBody, err
    }

    response := &SportResponse{
        ID:   sport.ID,
        Name: sport.Name,
    }

    return &utils.ResponseBody[SportResponse]{
        Body: response,
    }, nil
}

// Connect to DB and delete sport
func (u *SportService) DeleteSport(ctx context.Context, input *struct {
	Body DeleteSportRequest
}) (*utils.ResponseBody[DeleteSportResponse], error) {
	err := u.sportDB.DeleteSport(input.Body.ID)
	respBody := &utils.ResponseBody[DeleteSportResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &DeleteSportResponse{
		ID: input.Body.ID,
	}

	return &utils.ResponseBody[DeleteSportResponse]{
		Body: response,
	}, nil
}