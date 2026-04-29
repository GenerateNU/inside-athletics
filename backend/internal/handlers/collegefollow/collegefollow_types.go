package collegefollow

import "github.com/google/uuid"

// Given UserID, get all colleges that are followed
type GetCollegeFollowsByUserParams struct {
}

// Given college, get list of UserIDs that follow this college
type GetFollowingUsersByCollegeParams struct {
	CollegeID uuid.UUID `path:"college_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the College"`
}

type GetCollegeFollowsByUserResponse struct {
	UserID    uuid.UUID   `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
	CollegeIDs []uuid.UUID `json:"college_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"The colleges the given user follows"`
}

type GetFollowingUsersByCollegeResponse struct {
	CollegeID uuid.UUID   `json:"college_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the College"`
	UserIDs   []uuid.UUID `json:"user_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"All users that follow this college"`
}

type CreateCollegeFollowInput struct {
	Body CreateCollegeFollowBody
}

type CreateCollegeFollowBody struct {
	CollegeID uuid.UUID `json:"college_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the College"`
}

type CreateCollegeFollowResponse struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the college follow created"`
	CollegeID uuid.UUID `json:"college_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the College"`
	UserID    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

type DeleteCollegeFollowParams struct {
	CollegeID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the college to unfollow"`
}

type DeleteCollegeFollowResponse struct {
	Message string `json:"message" example:"College follow was deleted successfully" doc:"Message to display"`
}
