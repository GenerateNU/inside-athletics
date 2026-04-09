package college

import (
	"context"
	"fmt"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

// Contains business logic for colleges
type CollegeService struct {
	collegeDB *CollegeDB
}

// Convert empty string to nil pointer
func StringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Retrieves single college by ID
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
		Website:      college.Website,
		AcademicRank: college.AcademicRank,
		DivisionRank: college.DivisionRank,
		Logo:         StringPtrOrNil(college.Logo),
	}

	return &utils.ResponseBody[GetCollegeResponse]{
		Body: response,
	}, err
}

func (u *CollegeService) ListColleges(ctx context.Context, input *ListCollegesParams) (*utils.ResponseBody[ListCollegesResponse], error) {
	colleges, err := u.collegeDB.ListColleges(input.Limit, input.Offset)

	respBody := &utils.ResponseBody[ListCollegesResponse]{}
	if err != nil {
		return respBody, err
	}

	responseColleges := make([]GetCollegeResponse, 0, len(*colleges))
	for _, college := range *colleges {
		responseColleges = append(responseColleges, GetCollegeResponse{
			ID:           college.ID,
			Name:         college.Name,
			State:        college.State,
			City:         college.City,
			Website:      college.Website,
			AcademicRank: college.AcademicRank,
			DivisionRank: college.DivisionRank,
			Logo:         StringPtrOrNil(college.Logo),
		})
	}

	return &utils.ResponseBody[ListCollegesResponse]{
		Body: &ListCollegesResponse{
			Colleges: responseColleges,
			Total:    len(responseColleges),
		},
	}, nil
}

// Creates a single college
func (u *CollegeService) CreateCollege(ctx context.Context, input *CreateCollegeInput) (*utils.ResponseBody[CreateCollegeResponse], error) {
	college := &models.College{
		Name:         input.Body.Name,
		State:        input.Body.State,
		City:         input.Body.City,
		Website:      input.Body.Website,
		DivisionRank: input.Body.DivisionRank,
		AcademicRank: input.Body.AcademicRank,
	}
	if input.Body.Logo != nil {
		college.Logo = *input.Body.Logo
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
		Website:      createdCollege.Website,
		AcademicRank: createdCollege.AcademicRank,
		DivisionRank: createdCollege.DivisionRank,
		Logo:         StringPtrOrNil(createdCollege.Logo),
	}

	return &utils.ResponseBody[CreateCollegeResponse]{
		Body: response,
	}, err
}

// Updates fields inputted in request
func (u *CollegeService) UpdateCollege(ctx context.Context, input *UpdateCollegeInput) (*utils.ResponseBody[UpdateCollegeResponse], error) {
	id := input.ID

	college, err := u.collegeDB.UpdateCollege(id, &input.Body)

	respBody := &utils.ResponseBody[UpdateCollegeResponse]{}
	if err != nil {
		return respBody, err
	}

	response := &UpdateCollegeResponse{
		ID:           college.ID,
		Name:         college.Name,
		State:        college.State,
		City:         college.City,
		Website:      college.Website,
		AcademicRank: college.AcademicRank,
		DivisionRank: college.DivisionRank,
		Logo:         StringPtrOrNil(college.Logo),
	}

	return &utils.ResponseBody[UpdateCollegeResponse]{
		Body: response,
	}, err
}

// Removes a single college by ID
func (u *CollegeService) DeleteCollege(ctx context.Context, input *DeleteCollegeParams) (*utils.ResponseBody[DeleteCollegeResponse], error) {
	id := input.ID
	err := u.collegeDB.DeleteCollege(id)

	respBody := &utils.ResponseBody[DeleteCollegeResponse]{}
	if err != nil {
		return respBody, err
	}

	response := &DeleteCollegeResponse{
		Message: fmt.Sprintf("College %s deleted successfully", id.String()),
		ID:      id,
	}

	return &utils.ResponseBody[DeleteCollegeResponse]{
		Body: response,
	}, err
}
