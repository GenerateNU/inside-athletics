package types

type HealthResponse struct {
	Message string `json:"message" example:"Healthy!" doc:"Message to display"`
}

type EmptyInput struct{}
