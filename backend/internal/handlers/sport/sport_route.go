package sport

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	sportService := NewSportService(db)

	{
		grp := huma.NewGroup(api, "/api/v1/sport")
		huma.Post(grp, "/", sportService.CreateSport)                 // Create sport
		huma.Get(grp, "/by-name/{name}", sportService.GetSportByName) // Read sport by name
		huma.Get(grp, "/{id}", sportService.GetSportByID)             // Read sport by ID
		huma.Patch(grp, "/{id}", sportService.UpdateSport)            // Update sport
		huma.Delete(grp, "/{id}", sportService.DeleteSport)           // Delete sport
	}
	{
		grp := huma.NewGroup(api, "/api/v1/sports")
		huma.Get(grp, "/", sportService.GetAllSports) // Read sports
	}
}
