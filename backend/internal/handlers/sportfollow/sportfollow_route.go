package sportfollow

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {

	var sportFollowDB = &SportFollowDB{db: db}
	var sportFollowService = &SportFollowService{sportfollowDB: sportFollowDB}
	{
		grp := huma.NewGroup(api, "/api/v1/user/sport")
		huma.Post(grp, "/", sportFollowService.CreateSportFollow)
		huma.Get(grp, "/follows", sportFollowService.GetSportFollowsByUser)
		huma.Get(grp, "/{sport_id}/users", sportFollowService.GetFollowingUsersBySport)
		huma.Delete(grp, "/{id}", sportFollowService.DeleteSportFollow)
	}
}
