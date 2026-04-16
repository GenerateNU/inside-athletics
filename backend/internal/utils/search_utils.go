package utils

import (
	"fmt"

	"gorm.io/gorm"
)

func FuzzySearchByQueries(colName string, searchStr string) (string, string, string) {
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

func FuzzySearchForDB[ModelType any](searchStr string, limit int, db *gorm.DB, colName string, t ModelType) ([]ModelType, error) {
	var searchResults []ModelType

	selectQuery, whereQuery, orderQuery := FuzzySearchByQueries(colName, searchStr)

	if err := db.Select("*", selectQuery).Where(whereQuery).Order(orderQuery).Limit(limit).Model(&t).Scan(&searchResults).Error; err != nil {
		return []ModelType{}, err
	}

	return searchResults, nil
}

func FuzzySearchService[ModelType any, RespType any](input *SearchParam, modelType ModelType, respType RespType, colName string, db *gorm.DB, toResp func(*ModelType) *RespType) (*ResponseBody[SearchResults[*RespType]], error) {
	searchResults, err := FuzzySearchForDB(input.SearchStr, input.Limit, db, "name", modelType)
	respBody := ResponseBody[SearchResults[*RespType]]{}
	if err != nil {
		return HandleDBError(&respBody, err)
	}

	// parse search results into response type
	searchResponses := make([]*RespType, 0)
	for _, m := range searchResults {
		searchResponses = append(searchResponses, toResp(&m))
	}
	respBody.Body = &SearchResults[*RespType]{
		Results: searchResponses,
	}
	return &respBody, nil
}

type SearchParam struct {
	SearchStr string `query:"search_str" binding:"required" example:"Northeastern University" docs:"Search string to find colleges by"`
	Limit     int    `query:"limit" default:"20" example:"20" docs:"Max number of search results to return"`
}

type SearchResults[T any] struct {
	Results []T `json:"results" docs:"List of search results"`
}
