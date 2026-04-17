package utility

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	utilityService := &UtilityService{NewUtilityDB(db)}

	grp := huma.NewGroup(api, "/api/v1/utility")
	huma.Get(grp, "/access-check", utilityService.GetAccessCheck)
}
