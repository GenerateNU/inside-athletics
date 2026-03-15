package sportfollow

import "github.com/google/uuid"

// Given UserID, get all sports that are followed
type GetSportFollowsByUserParams struct {
	UserID uuid.UUID `path:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

// Given sport, get list of UserIDs that follow this sport
type GetFollowingUsersBySportParams struct {
	SportID uuid.UUID `path:"sport_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Sport"`
}

type GetSportFollowsByUserResponse struct {
	UserID   uuid.UUID   `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
	SportIDs []uuid.UUID `json:"sport_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"The sports the given user follows"`
}

type GetFollowingUsersBySportResponse struct {
	SportID uuid.UUID   `json:"sport_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Sport"`
	UserIDs []uuid.UUID `json:"user_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"All users that follow this sport"`
}

type CreateSportFollowInput struct {
	Body CreateSportFollowBody
}

type CreateSportFollowBody struct {
	SportID uuid.UUID `json:"sport_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Sport"`
	UserID  uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

type CreateSportFollowResponse struct {
	ID      uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the sport follow created"`
	SportID uuid.UUID `json:"sport_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Sport"`
	UserID  uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

type DeleteSportFollowParams struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the sport follow to be deleted"`
}

type DeleteSportFollowResponse struct {
	Message string `json:"message" example:"Sport follow was deleted successfully" doc:"Message to display"`
}
