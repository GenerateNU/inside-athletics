package goat // take note that the package matches the directory name

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	goatDB := &GoatDB{db: db}                  // create instance of the db!
	goatService := GoatService{goatDB: goatDB} // define an instance of the goat service
	{
		grp := huma.NewGroup(api, "api/v1/goat")    // this creates a goat route group!
		huma.Get(grp, "/", goatService.Ping)        // pass the group into the huma get method
		huma.Get(grp, "/{id}", goatService.GetGoat) // NEW ENDPOINT!
	}
}
