package premiumpost

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var premiumPostDB = &PremiumPostDB{db: db}
	var premiumPostService = &PremiumPostService{premiumPostDB}
	{
		grp := huma.NewGroup(api, "/api/v1/post/premium")
		huma.Post(grp, "/", premiumPostService.CreatePremiumPost) // Create post
	}
	{
		grp := huma.NewGroup(api, "/api/v1/posts/premium")
		huma.Get(grp, "/", premiumPostService.GetAllPremiumPosts)
		huma.Get(grp, "/by-author/{author_id}", premiumPostService.GetPremiumPostsByAuthorID)
		huma.Get(grp, "/by-sport/{sport_id}", premiumPostService.GetPremiumPostsBySportID)
		huma.Get(grp, "/by-college/{college_id}", premiumPostService.GetPremiumPostsByCollegeID)
		huma.Get(grp, "/by-tag/{tag_id}", premiumPostService.GetPremiumPostsByTagID)
	}
}
