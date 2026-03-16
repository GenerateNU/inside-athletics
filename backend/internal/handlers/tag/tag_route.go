package tag

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"inside-athletics/internal/handlers/tagpost"
)

func Route(api huma.API, db *gorm.DB) {
	var tagDB = &TagDB{db} // create object storing all database level functions for user
	var tagPostDB = tagpost.NewTagPostDB(db)
	var tagService = &TagService{tagDB, tagPostDB} // create object with user functionality
	{
		grp := huma.NewGroup(api, "/api/v1/tag")
		huma.Post(grp, "/", tagService.CreateTag)
		huma.Get(grp, "/name/{name}", tagService.GetTagByName)
		huma.Get(grp, "/{id}", tagService.GetTagById)
		huma.Get(grp, "/{tag_id}/posts", tagService.GetPostsByTag)
		huma.Patch(grp, "/{id}", tagService.UpdateTag)
		huma.Delete(grp, "/{id}", tagService.DeleteTag)
	}
}
