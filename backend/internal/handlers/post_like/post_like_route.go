package post_like

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var postLikeDB = &PostLikeDB{db: db}
	var postLikeService = &PostLikeService{postLikeDB}
	{
		grp := huma.NewGroup(api, "/api/v1/post/like")
		huma.Post(grp, "", postLikeService.CreatePostLike)                // Create like
		huma.Get(grp, "/{id}", postLikeService.GetPostLike)                // Get like by ID
		huma.Delete(grp, "/{id}", postLikeService.DeletePostLike)          // Delete like
		huma.Get(grp, "/{post_id}/likes", postLikeService.GetPostLikeInfo) // Like count and whether user liked
	}
}
