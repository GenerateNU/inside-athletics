

type GetUserParams struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}

type GetUserResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}