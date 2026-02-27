package tagfollow

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {

	var tagFollowDB = &TagFollowDB{db: db}
	var tagFollowService = &TagFollowService{tagFollowDB}
	{
		grp := huma.NewGroup(api, "/api/v1/user/tag")
		huma.Post(grp, "/", tagFollowService.CreateTagFollow)
		huma.Get(grp, "/{user_id}/follows", tagFollowService.GetTagFollowsByUser)
		huma.Get(grp, "/{tag_id}/users", tagFollowService.GetFollowingUsersByTag)
		huma.Delete(grp, "/{id}", tagFollowService.DeleteTagFollow)
	}
}
