package models

type HealthModel struct {
	Id   int8   `json:"id" example:"1" doc:"ID of this health data"`
	Word string `json:"word" example:"erm" doc:"test word for example data"`
}
