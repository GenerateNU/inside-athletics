package collegefollow

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {

	var collegeFollowDB = &CollegeFollowDB{db: db}
	var collegeFollowService = &CollegeFollowService{collegefollowDB: collegeFollowDB}
	{
		grp := huma.NewGroup(api, "/api/v1/user/college")
		huma.Post(grp, "/", collegeFollowService.CreateCollegeFollow)
		huma.Get(grp, "/{user_id}/follows", collegeFollowService.GetCollegeFollowsByUser)
		huma.Get(grp, "/{college_id}/users", collegeFollowService.GetFollowingUsersByCollege)
		huma.Delete(grp, "/{id}", collegeFollowService.DeleteCollegeFollow)
	}
}
