package college

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

type CollegeService struct {
	collegeDB *CollegeDB
}

func (u *CollegeService) GetCollege(ctx context.Context, input *GetCollegeParams) (*utils.ResponseBody[GetCollegeResponse], error) {
	id := input.ID
	college, err := u.collegeDB.GetCollege(id)
	respBody := &utils.ResponseBody[GetCollegeResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &GetCollegeResponse{
		ID:           college.ID,
		Name:         college.Name,
		State:        college.State,
		City:         college.City,
		Website:      &college.Website,
		AcademicRank: college.AcademicRank,
		DivisionRank: college.DivisionRank,
		Logo:         &college.Logo,
	}

	return &utils.ResponseBody[GetCollegeResponse]{
		Body: response,
	}, err
}

func (u *CollegeService) CreateCollege(ctx context.Context, input *CreateCollegeRequest) (*utils.ResponseBody[CreateCollegeResponse], error) {
	college := &models.College{
		Name:         input.Name,
		State:        input.State,
		City:         input.City,
		DivisionRank: input.DivisionRank,
	}

	createdCollege, err := u.collegeDB.CreateCollege(college)
	respBody := &utils.ResponseBody[CreateCollegeResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &CreateCollegeResponse{
		ID:           createdCollege.ID,
		Name:         createdCollege.Name,
		State:        createdCollege.State,
		City:         createdCollege.City,
		Website:      &createdCollege.Website,
		AcademicRank: createdCollege.AcademicRank,
		DivisionRank: createdCollege.DivisionRank,
		Logo:         &createdCollege.Logo,
	}

	return &utils.ResponseBody[CreateCollegeResponse]{
		Body: response,
	}, err
}

func (u *CollegeService) UpdateCollege(ctx context.Context, input *UpdateCollegeInput) (*utils.ResponseBody[UpdateCollegeResponse], error) {
	id := input.ID
	updates := make(map[string]interface{})

	college, err := u.collegeDB.UpdateCollege(id, updates)
	respBody := &utils.ResponseBody[UpdateCollegeResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &UpdateCollegeResponse{
		ID:           college.ID,
		Name:         college.Name,
		State:        college.State,
		City:         college.City,
		Website:      &college.Website,
		AcademicRank: college.AcademicRank,
		DivisionRank: college.DivisionRank,
		Logo:         &college.Logo,
	}

	return &utils.ResponseBody[UpdateCollegeResponse]{
		Body: response,
	}, err
}

func (u *CollegeService) DeleteCollege(ctx context.Context, input *DeleteCollegeParams) (*utils.ResponseBody[DeleteCollegeResponse], error) {
	id := input.ID
	err := u.collegeDB.DeleteCollege(id)
	respBody := &utils.ResponseBody[DeleteCollegeResponse]{}

	if err != nil {
		return respBody, err
	}

	response := &DeleteCollegeResponse{
		Message: "College deleted successfully",
	}

	return &utils.ResponseBody[DeleteCollegeResponse]{
		Body: response,
	}, err
}
