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
		huma.Post(grp, "/", premiumPostService.CreatePremiumPost)                 // Create post
	}
}