package utils

import (
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type RouteFN func(api huma.API, db *gorm.DB)

func FuzzySearchBy(colName string, searchStr string) (string, string, string) {
	// Create query strings for Select, Where, and Order By arguments
	selectQuery := fmt.Sprintf(
		"word_similarity('%s', %s) as similarity",
		searchStr,
		colName,
	)
	whereQuery := fmt.Sprintf(
		"word_similarity('%s', %s) >= show_limit()",
		searchStr,
		colName,
	)

	return selectQuery, whereQuery, "similarity DESC"
}
