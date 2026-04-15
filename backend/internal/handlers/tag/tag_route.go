package tag

import (
	"inside-athletics/internal/handlers/tagpost"
	"inside-athletics/internal/s3"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB, s3Svc *s3.Service) {
	var tagDB = &TagDB{db} // create object storing all database level functions for user
	var tagPostDB = tagpost.NewTagPostDB(db)
	var tagService = &TagService{tagDB, tagPostDB, s3Svc} // create object with user functionality
	{
		grp := huma.NewGroup(api, "/api/v1/tag")
		huma.Post(grp, "/", tagService.CreateTag)
		huma.Get(grp, "/name/{name}", tagService.GetTagByName)
		huma.Get(grp, "/{id}", tagService.GetTagById)
		huma.Get(grp, "/{tag_id}/posts", tagService.GetPostsByTag)
		huma.Get(grp, "/type/{type}", tagService.GetTagsByType)
		huma.Patch(grp, "/{id}", tagService.UpdateTag)
		huma.Delete(grp, "/{id}", tagService.DeleteTag)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/tags")
		huma.Get(grp, "/search", tagService.FuzzySearchFor)
	}
}
