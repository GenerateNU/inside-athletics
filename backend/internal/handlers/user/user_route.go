package user

import (
	"inside-athletics/internal/handlers/role"
	"inside-athletics/internal/s3"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

/*
*
Groups together all of the User routes. Huma is a wrapper here that automatically does a few things:

 1. creates OpenAPI docs
 2. maps the response to the correct response type (if no error 200, 201, etc.) if error it will use the Huma
    error status code
*/
func Route(api huma.API, db *gorm.DB, s3Service *s3.Service) {
	var userDB = NewUserDB(db)
	var roleDB = role.NewRoleDB(db) // create object storing all database level functions for user
	var userService = &UserService{userDB: userDB, roleDB: roleDB, s3: s3Service}
	{
		grp := huma.NewGroup(api, "/api/v1/user")
		huma.Get(grp, "/current", userService.GetCurrentUser)
		huma.Post(grp, "", userService.CreateUser)
		huma.Get(grp, "/{id}", userService.GetUser)
		huma.Patch(grp, "", userService.UpdateUser)
		huma.Delete(grp, "/{id}", userService.DeleteUser)
		huma.Post(grp, "/{id}/roles", userService.AssignRole)
	}
}
