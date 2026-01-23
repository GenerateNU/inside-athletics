package sport

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	sportDB := &SportDB{db: db}
	sportService := SportService{sportDB: sportDB}
	{
		grp := huma.NewGroup(api, "api/v1/sport")
		huma.Post(grp,		"/", 		sportService.CreateSport)       // C reate sport
		huma.Get(grp, 		"/{id}", 	sportService.GetSportById) 		// R ead sport
		huma.Get(grp, 		"/{name}", 	sportService.GetSportByName)	// R ead sport
		huma.Patch(grp, 	"/{id}", 	sportService.UpdateSport)  		// U pdate sport
		huma.Delete(grp, 	"/{id}", 	sportService.DeleteSport) 		// D elete sport
	}
	{
		grp := huma.NewGroup(api, "api/v1/sports")
		huma.Get(grp, 		"/", sportService.GetAllSports) 			// R ead sports
	}
}
