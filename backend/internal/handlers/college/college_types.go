package college

import (
	"github.com/google/uuid"

	"inside-athletics/internal/models"
)

type GetCollegeParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID of the college"`
}

type ListCollegesParams struct {
	Limit  int `query:"limit" default:"200" example:"200" doc:"Maximum number of colleges to return"`
	Offset int `query:"offset" default:"0" example:"0" doc:"Number of colleges to skip"`
}

type GetCollegeResponse struct {
	ID           uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the college"`
	Name         string          `json:"name" example:"Northeastern University" doc:"Name of the college"`
	State        string          `json:"state" example:"Massachusetts" doc:"State of the college"`
	City         string          `json:"city"  example:"Boston" doc:"City of the college"`
	Website      string          `json:"website" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16          `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank models.Division `json:"division_rank" enum:"1,2,3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string         `json:"logo" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type ListCollegesResponse struct {
	Colleges []GetCollegeResponse `json:"colleges" doc:"List of colleges"`
	Total    int                  `json:"total" doc:"Total number of colleges returned"`
}

type CreateCollegeRequest struct {
	Name         string          `json:"name" required:"true" minLength:"1" maxLength:"200" example:"Northeastern University" doc:"Name of the college"`
	State        string          `json:"state" required:"true" minLength:"1" maxLength:"100" example:"Massachusetts" doc:"State of the college"`
	City         string          `json:"city" required:"true" minLength:"1" maxLength:"100" example:"Boston" doc:"City of the college"`
	Website      string          `json:"website" required:"true" minLength:"1" maxLength:"500" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16          `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank models.Division `json:"division_rank" required:"true" enum:"1,2,3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string         `json:"logo" maxLength:"500" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type CreateCollegeResponse struct {
	ID           uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the college"`
	Name         string          `json:"name" example:"Northeastern University" doc:"Name of the college"`
	State        string          `json:"state" example:"Massachusetts" doc:"State of the college"`
	City         string          `json:"city" example:"Boston" doc:"City of the college"`
	Website      string          `json:"website" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16          `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank models.Division `json:"division_rank" example:"1" enum:"1,2,3" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string         `json:"logo" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type CreateCollegeInput struct {
	Body CreateCollegeRequest
}

type UpdateCollegeRequest struct {
	Name         *string          `json:"name" maxLength:"200" example:"Northeastern University" doc:"Name of the college"`
	State        *string          `json:"state" maxLength:"100" example:"Massachusetts" doc:"State of the college"`
	City         *string          `json:"city" maxLength:"100" example:"Boston" doc:"City of the college"`
	Website      *string          `json:"website" maxLength:"500" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16           `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank *models.Division `json:"division_rank" enum:"1,2,3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string          `json:"logo" maxLength:"500" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

// Combined input for Update (path params + body)
type UpdateCollegeInput struct {
	ID   uuid.UUID `path:"id" example:"1" doc:"ID of the college"`
	Body UpdateCollegeRequest
}

type UpdateCollegeResponse struct {
	ID           uuid.UUID       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the college"`
	Name         string          `json:"name" required:"true" example:"Northeastern University" doc:"Name of the college"`
	State        string          `json:"state" required:"true" example:"Massachusetts" doc:"State of the college"`
	City         string          `json:"city" required:"true" example:"Boston" doc:"City of the college"`
	Website      string          `json:"website" required:"true" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16          `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank models.Division `json:"division_rank" required:"true" enum:"1,2,3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string         `json:"logo" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type DeleteCollegeParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID of the college"`
}

type DeleteCollegeResponse struct {
	Message string    `json:"message" example:"College deleted successfully" doc:"Success message"`
	ID      uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the deleted college"`
}
