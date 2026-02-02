package post

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	postService := NewPostService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/post")
		huma.Post(grp, "/", postService.CreatePost)                            // Create post
		huma.Get(grp, "/by-title/{title}", postService.GetPostByTitle)         // Read post by title
		huma.Get(grp, "/{id}", postService.GetPostByID)                        // Read post by ID
		huma.Get(grp, "/by-author/{author_id}", postService.GetPostByAuthorID) // Read post by author id
		huma.Get(grp, "by-sport/{sport_id}", postService.GetPostBySportID)     // Read post by sport id
		huma.Patch(grp, "/{id}", postService.UpdatePost)                       // Update post
		huma.Delete(grp, "/{id}", postService.DeletePost)                      // Delete post
	}
	{
		grp := huma.NewGroup(api, "/api/v1/posts")
		huma.Get(grp, "/", postService.GetAllPosts) // Read all posts
	}
}
