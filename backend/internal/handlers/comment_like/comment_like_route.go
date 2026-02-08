package like

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var commentLikeDB = &CommentLikeDB{db}               // create object storing all database level functions for like
	var commentLikeService = &CommentLikeService{commentLikeDB} // create object with like functionality
	{
		grp := huma.NewGroup(api, "/api/v1/user")
		huma.Post(grp, "/", commentLikeService.CreateCommentLike)
		huma.Get(grp, "/{id}", commentLikeService.GetCommentLike)
		huma.Delete(grp, "/{id}", commentLikeService.DeleteCommentLike)
		huma.Get(grp, "/comment/{comment_id}/like-count", commentLikeService.GetLikeCount)
		huma.Get(grp, "/comment/{comment_id}/check-like", commentLikeService.CheckUserLikedComment)
	}
}