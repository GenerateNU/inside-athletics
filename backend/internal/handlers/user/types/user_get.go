package types

type GetUserParams struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}

type GetUserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
}