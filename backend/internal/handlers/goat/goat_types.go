package goat

type GoatResponse struct {
   Id int8 `json:"id" example:"1" doc:"id of this goat"`
   Name string `json:"name" example:"Suli" doc:"name of this goat"`
   Age string `json:"age" example:"67" doc:"age of this goat"`
}

type GetGoatParams struct {
	Id int8 `path:"id" example:"1" doc:"id of this goat"`
}