package healthRouteTypes

import "inside-athletics/internal/models"

type GetHealthInput struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}

type GetHealthOutput struct {
	Body *models.HealthModel
}
