package user

type GetUserParams struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}

type GetUserResponse struct {
	ID   uint   `json:"id" example:"1" doc:"ID of the user"`
	Name string `json:"name" example:"Joe" doc:"Name of the user"`
}
