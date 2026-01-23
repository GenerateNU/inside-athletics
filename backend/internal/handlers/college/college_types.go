package college

import "github.com/google/uuid"

type GetCollegeParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID of the college"`
}

type GetCollegeResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the college"`
	Name         string    `json:"name" example:"Northeastern University" doc:"Name of the college"`
	State        string    `json:"state" example:"Massachusetts" doc:"State of the college"`
	City         string    `json:"city" example:"Boston" doc:"City of the college"`
	Website      *string   `json:"website" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16    `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank int8      `json:"division_rank" minimum:"1" maximum:"3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string   `json:"logo" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type CreateCollegeRequest struct {
	Name         string  `json:"name" maxLength:"200" example:"Northeastern University" doc:"Name of the college"`
	State        string  `json:"state" maxLength:"100" example:"Massachusetts" doc:"State of the college"`
	City         string  `json:"city" maxLength:"100" example:"Boston" doc:"City of the college"`
	Website      *string `json:"website" maxLength:"500" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16  `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank int8    `json:"division_rank" minimum:"1" maximum:"3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string `json:"logo" maxLength:"500" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type CreateCollegeResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" doc:"ID of the college"`
	Name         string    `json:"name" example:"Northeastern University" doc:"Name of the college"`
	State        string    `json:"state" example:"Massachusetts" doc:"State of the college"`
	City         string    `json:"city" example:"Boston" doc:"City of the college"`
	Website      *string   `json:"website" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16    `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank int8      `json:"division_rank" minimum:"1" maximum:"3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string   `json:"logo" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

type UpdateCollegeParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID of the college"`
}

type UpdateCollegeRequest struct {
	Name         *string `json:"name" maxLength:"200" example:"Northeastern University" doc:"Name of the college"`
	State        *string `json:"state" maxLength:"100" example:"Massachusetts" doc:"State of the college"`
	City         *string `json:"city" maxLength:"100" example:"Boston" doc:"City of the college"`
	Website      *string `json:"website" maxLength:"500" example:"https://www.northeastern.edu" doc:"Website of the college"`
	AcademicRank *int16  `json:"academic_rank" example:"53" doc:"Academic rank of the college"`
	DivisionRank *int8   `json:"division_rank" minimum:"1" maximum:"3" example:"1" doc:"NCAA division (1, 2, or 3)"`
	Logo         *string `json:"logo" maxLength:"500" example:"https://example.com/logo.png" doc:"Logo of the college"`
}

// Combined input for Update (path params + body)
type UpdateCollegeInput struct {
	UpdateCollegeParams
	UpdateCollegeRequest
}

type DeleteCollegeParams struct {
	ID uuid.UUID `path:"id" example:"1" doc:"ID of the college"`
}
