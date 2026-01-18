package goat

import (
   "github.com/danielgtaylor/huma/v2"
   "gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	goatDB := &GoatDB{db: db}
	goatService := GoatService{goatDB: goatDB}
	{
		grp := huma.NewGroup(api, "api/v1/goat")
		huma.Get(grp, "/", goatService.Ping)
		huma.Get(grp, "/{id}", goatService.GetGoat) // ADD THE FUNCTIONALITY
	}
}
