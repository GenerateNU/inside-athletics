package user

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var userDB *UserDB = &UserDB{db}              // create object storing all database level functions for user
	var UserService *UserService = &UserService{healthDB} // create object with user functionality
	{
		grp := huma.NewGroup(api, "/api/v1/user")
		huma.Get(grp, "/{name}", healthService.GetUser)
	}
}
