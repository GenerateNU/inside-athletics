package tagpost

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var tagpostDB = NewTagPostDB(db)
	var tagService = &TagPostService{tagpostDB}
	{
		grp := huma.NewGroup(api, "/api/v1/post/tag")
		huma.Post(grp, "/", tagService.CreateTagPost)
		huma.Get(grp, "/{id}", tagService.GetTagPostById)
		huma.Patch(grp, "/{id}", tagService.UpdateTagPost)
		huma.Delete(grp, "/{id}", tagService.DeleteTagPost)
	}
}
