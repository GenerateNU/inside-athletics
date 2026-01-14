package health

type HealthResponse struct {
	Message string `json:"message" example:"Healthy!" doc:"Message to display"`
}
