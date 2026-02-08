package like

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var postLikeDB = &PostLikeDB{db}
	var postLikeService = &PostLikeService{postLikeDB}
	{
		grp := huma.NewGroup(api, "/api/v1/user")
		huma.Post(grp, "/", postLikeService.CreatePostLike)
		huma.Get(grp, "/{id}", postLikeService.GetPostLike)
		huma.Delete(grp, "/{id}", postLikeService.DeletePostLike)
		huma.Get(grp, "/post/{post_id}/like-count", postLikeService.GetLikeCount)
		huma.Get(grp, "/post/{post_id}/check-like", postLikeService.CheckUserLikedPost)
	}
}
