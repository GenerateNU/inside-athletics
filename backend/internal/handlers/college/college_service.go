package college

import (
	"context"
	"inside-athletics/internal/utils"
)

type CollegeService struct {
	collegeDB *CollegeDB
}

/*
*
Huma automatically validates the params based on which type you have passed it based on the struct tags (so nice!!)
IF there is an error it will automatically send the correct response back with the right status and message about the validation errors

Here we are mapping to a GetCollegeResponse so that we can control what the return type is. It's important to make seperate
return types so that we can control what information we are sending back instead of just the entire model
*/
func (u *CollegeService) GetCollege(ctx context.Context, input *GetCollegeParams) (*utils.ResponseBody[GetCollegeResponse], error) {
	id := input.ID
	college, err := u.collegeDB.GetCollege(id)
	respBody := &utils.ResponseBody[GetCollegeResponse]{}

	if err != nil {
		return respBody, err
	}

	// mapping to correct response type
	// we do this so we can control what values are
	// returned by the API
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
