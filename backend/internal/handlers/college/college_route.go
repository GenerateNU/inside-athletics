package college

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

/*
*
Groups together all of the College routes. Huma is a wrapper here that automatically does a few things:

 1. creates OpenAPI docs
 2. maps the response to the correct response type (if no error 200, 201, etc.) if error it will use the Huma
    error status code
*/
func Route(api huma.API, db *gorm.DB) {
	var collegeDB = &CollegeDB{db}                  // create object storing all database level functions for college
	var collegeService = &CollegeService{collegeDB} // create object with college functionality
	{
		grp := huma.NewGroup(api, "/api/v1/college")
		huma.Get(grp, "/{id}", collegeService.GetCollege)
		huma.Post(grp, "", collegeService.CreateCollege)
		huma.Put(grp, "/{id}", collegeService.UpdateCollege)
		huma.Delete(grp, "/{id}", collegeService.DeleteCollege)
	}
}
