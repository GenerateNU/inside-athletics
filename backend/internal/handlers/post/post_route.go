package post

import (
	"inside-athletics/internal/handlers/user"
	"inside-athletics/internal/s3"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB, s3Svc *s3.Service) {
	userDB := user.NewUserDB(db)
	postService := NewPostService(db, userDB, s3Svc)
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
		huma.Get(grp, "/filter", postService.FilterPosts)                      // Filter for posts based on college, sport, and tags
	}
}
