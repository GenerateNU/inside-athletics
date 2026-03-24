package premiumpost

import (
	"context"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
)

// --- reminders
// add permission to only allow moderators and admins for access to get, create, update, delete
// add permission so that only paid users can view the posts

type PremiumPostService struct {
	premiumPostDB *PremiumPostDB
}

func (s *PremiumPostService) CreatePremiumPost(ctx context.Context, input *struct{ Body CreatePremiumPostParams }) (*utils.ResponseBody[CreatePremiumPostResponse], error) {
	id, err := utils.GetCurrentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if len(input.Body.Tags) == 0 && input.Body.SportID == nil && input.Body.CollegeID == nil {
		return nil, huma.Error400BadRequest("Need to have at least a single tag on a post")
	}
	premiumPost := &models.PremiumPost{
		AuthorID:       id,
		SportID:        input.Body.SportID,
		CollegeID:      input.Body.CollegeID,
		Title:          input.Body.Title,
		Content:        input.Body.Content,
		AttachmentKey:  input.Body.AttachmentKey,
		AttachmentType: input.Body.AttachmentType,
	}

	createdPost, err := utils.HandleDBError(
		s.premiumPostDB.CreatePremiumPost(premiumPost),
	)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[CreatePremiumPostResponse]{
		Body: ToCreatePremiumPostResponse(createdPost, id),
	}, nil
}
