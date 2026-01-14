package goat

type GoatResponse struct {
	Id   int8   `json:"id" example:"1" doc:"id of this goat"`
	Name string `json:"name" example:"Suli" doc:"Name of this goat"`
	Age  int8   `json:"age" example:"67" doc:"Age of this goat"`
}
