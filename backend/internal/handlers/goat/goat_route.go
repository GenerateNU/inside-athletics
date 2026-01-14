package goat

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	goat_service := GoatService{} // define an instance of the goat service
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", goat_service.Ping) // functionality of the endpoint!
	}
}
