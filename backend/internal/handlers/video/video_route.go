package video

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	videoService := NewVideoService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/video")
		huma.Post(grp, "/", videoService.CreateVideo)       // Add video
		huma.Get(grp, "/{id}", videoService.GetVideo)   // Get video by id
		huma.Delete(grp, "/{id}", videoService.DeleteVideo) // Delete video
	}
}
