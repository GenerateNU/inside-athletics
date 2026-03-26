package post

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	postService := NewPostService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/post")
		huma.Post(grp, "/", postService.CreatePost)       // Create post
		huma.Get(grp, "/{id}", postService.GetPostByID)   // Read post by ID
		huma.Patch(grp, "/{id}", postService.UpdatePost)  // Update post
		huma.Delete(grp, "/{id}", postService.DeletePost) // Delete post
	}
	{
		grp := huma.NewGroup(api, "/api/v1/posts")
		huma.Get(grp, "/", postService.GetAllPosts)                            // Read all posts
		huma.Get(grp, "/popular", postService.GetPopularPosts)                 // Read popular posts
		huma.Get(grp, "/by-sport/{sport_id}", postService.GetPostBySportID)    // Read posts by sport id
		huma.Get(grp, "/by-author/{author_id}", postService.GetPostByAuthorID) // Read posts by author id
		huma.Get(grp, "/search", postService.FuzzySearchForPost)               // Find all posts based on title for given search string
	}
}
