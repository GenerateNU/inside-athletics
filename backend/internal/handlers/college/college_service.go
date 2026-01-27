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

// Creates a map of non-nil fields to update
func buildUpdateMap(body *UpdateCollegeRequest) map[string]interface{} {
	updates := make(map[string]interface{})

	if body.Name != nil {
		updates["name"] = *body.Name
	}
	if body.State != nil {
		updates["state"] = *body.State
	}
	if body.City != nil {
		updates["city"] = *body.City
	}
	if body.Website != nil {
		updates["website"] = *body.Website
	}
	if body.AcademicRank != nil {
		updates["academic_rank"] = *body.AcademicRank
	}
	if body.DivisionRank != nil {
		updates["division_rank"] = *body.DivisionRank
	}
	if body.Logo != nil {
		updates["logo"] = *body.Logo
	}

	return updates
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
	updates := buildUpdateMap(&input.Body)

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
