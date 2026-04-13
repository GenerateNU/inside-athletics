package media

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	mediaService := NewMediaService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/media")
		huma.Post(grp, "/", mediaService.CreateMedia)       // Add media
		huma.Get(grp, "/{id}", mediaService.GetMedia)       // Get media by id
		huma.Delete(grp, "/{id}", mediaService.DeleteMedia) // Delete media
	}
}
