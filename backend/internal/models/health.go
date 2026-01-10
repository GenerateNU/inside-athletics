package models

type HealthModel struct {
	Id   int8   `json:"id" example:"1" doc:"ID of this health data"`
	Name string `json:"name" example:"Joe" doc:"test name for example data"`
}
