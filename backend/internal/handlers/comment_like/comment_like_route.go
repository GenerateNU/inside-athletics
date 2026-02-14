package comment_like

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var commentLikeDB = &CommentLikeDB{db: db}
	var commentLikeService = &CommentLikeService{commentLikeDB}
	{
		grp := huma.NewGroup(api, "/api/v1/comment-like")
		huma.Post(grp, "/", commentLikeService.CreateCommentLike)                    // Create like
		huma.Get(grp, "/{id}", commentLikeService.GetCommentLike)                    // Get like by ID
		huma.Delete(grp, "/{id}", commentLikeService.DeleteCommentLike)               // Delete like
		huma.Get(grp, "/comment/{comment_id}/likes", commentLikeService.GetCommentLikeInfo) // Like count and whether user liked
	}
}