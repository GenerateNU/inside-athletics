package comment

import (
	"inside-athletics/internal/s3"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB, s3Svc *s3.Service) {
	commentDB := NewCommentDB(db)
	commentService := &CommentService{
		commentDB: commentDB,
		s3:        s3Svc,
	}
	{
		grp := huma.NewGroup(api, "/api/v1/comment")
		huma.Post(grp, "/", commentService.CreateComment)         // Create comment
		huma.Get(grp, "/{id}", commentService.GetComment)         // Get comment by ID
		huma.Get(grp, "/{id}/replies", commentService.GetReplies) // Get replies to comment
		huma.Patch(grp, "/{id}", commentService.UpdateComment)    // Update comment
		huma.Delete(grp, "/{id}", commentService.DeleteComment)   // Delete comment
	}
	{
		grp := huma.NewGroup(api, "/api/v1/post")
		huma.Get(grp, "/{post_id}/comments", commentService.GetCommentsByPost) // List comments by post
	}
}
