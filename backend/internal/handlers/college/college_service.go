package college

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"
)

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
		Website:      StringPtrOrNil(college.Website),
		AcademicRank: college.AcademicRank,
		DivisionRank: college.DivisionRank,
		Logo:         StringPtrOrNil(college.Logo),
	}

	return &utils.ResponseBody[GetCollegeResponse]{
		Body: response,
	}, err
}

func (u *CollegeService) CreateCollege(ctx context.Context, input *CreateCollegeInput) (*utils.ResponseBody[CreateCollegeResponse], error) {
	college := &models.College{
		Name:         input.Body.Name,
		State:        input.Body.State,
		City:         input.Body.City,
		DivisionRank: input.Body.DivisionRank,
	}

	// If provided in request, use them, otherwise leave as empty string/nil
	if input.Body.Website != nil {
		college.Website = *input.Body.Website
	}
	college.AcademicRank = input.Body.AcademicRank
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
		Website:      StringPtrOrNil(createdCollege.Website),
		AcademicRank: createdCollege.AcademicRank,
		DivisionRank: createdCollege.DivisionRank,
		Logo:         StringPtrOrNil(createdCollege.Logo),
	}

	return &utils.ResponseBody[CreateCollegeResponse]{
		Body: response,
	}, err
}

func (u *CollegeService) UpdateCollege(ctx context.Context, input *UpdateCollegeInput) (*utils.ResponseBody[UpdateCollegeResponse], error) {
	id := input.ID
	updates := make(map[string]interface{})

	// Only include fields that are provided
	if input.Body.Name != nil {
		updates["name"] = *input.Body.Name
	}
	if input.Body.State != nil {
		updates["state"] = *input.Body.State
	}
	if input.Body.City != nil {
		updates["city"] = *input.Body.City
	}
	if input.Body.Website != nil {
		updates["website"] = *input.Body.Website
	}
	if input.Body.AcademicRank != nil {
		updates["academic_rank"] = *input.Body.AcademicRank
	}
	if input.Body.DivisionRank != nil {
		updates["division_rank"] = *input.Body.DivisionRank
	}
	if input.Body.Logo != nil {
		updates["logo"] = *input.Body.Logo
	}

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
		Website:      StringPtrOrNil(college.Website),
		AcademicRank: college.AcademicRank,
		DivisionRank: college.DivisionRank,
		Logo:         StringPtrOrNil(college.Logo),
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
