package utils

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type RouteFN func(api huma.API, db *gorm.DB)
