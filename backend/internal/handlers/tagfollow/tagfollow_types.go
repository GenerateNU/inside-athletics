package tagfollow

import "github.com/google/uuid"

// Given UserID, get all tags that are followed
type GetTagFollowsByUserParams struct {
	UserID uuid.UUID `path:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

// Given tag, get list of UserIDs that follow this tag
type GetFollowingUsersByTagParams struct {
	TagID uuid.UUID `path:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
}

type GetTagFollowsByUserResponse struct {
	UserID uuid.UUID   `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
	TagIDs []uuid.UUID `json:"tag_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"The tags the given user follows"`
}

type GetFollowingUsersByTagResponse struct {
	TagID   uuid.UUID   `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
	UserIDs []uuid.UUID `json:"user_ids" example:"[\"123e4567-e89b-12d3-a456-426614174000\",\"123e4567-e89b-12d3-a456-426614174001\"]" doc:"All users that follow this tag"`
}

type CreateTagFollowInput struct {
	Body CreateTagFollowBody
}

type CreateTagFollowBody struct {
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

type CreateTagFollowResponse struct {
	ID     uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tag follow created"`
	TagID  uuid.UUID `json:"tag_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the Tag"`
	UserID uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the user"`
}

type DeleteTagFollowParams struct {
	ID uuid.UUID `path:"id" example:"123e4567-e89b-12d3-a456-426614174000" doc:"ID of the tag follow to be deleted"`
}

type DeleteTagFollowResponse struct {
	Message string `json:"message" example:"Tag follow was deleted successfully" doc:"Message to display"`
}
