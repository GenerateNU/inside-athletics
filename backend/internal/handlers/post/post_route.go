package post

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	postService := NewPostService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/post")
		huma.Post(grp, "/", postService.CreatePost)                 // Create post
		huma.Get(grp, "/{id}", postService.GetPostByID)             // Read post by ID
		huma.Get(grp, "/{post_id}/tags", postService.GetTagsByPost) // Read tags by post id
		huma.Patch(grp, "/{id}", postService.UpdatePost)            // Update post
		huma.Delete(grp, "/{id}", postService.DeletePost)           // Delete post
	}
	{
		grp := huma.NewGroup(api, "/api/v1/posts")
		huma.Get(grp, "/", postService.GetAllPosts)                            // Read all posts
		huma.Get(grp, "/by-sport/{sport_id}", postService.GetPostBySportID)    // Read posts by sport id
		huma.Get(grp, "/by-author/{author_id}", postService.GetPostByAuthorID) // Read posts by author id
	}
}
