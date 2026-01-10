package healthRouteParams

type GetHealthParams struct {
	Name string `path:"name" maxLength:"30" example:"Joe" doc:"Name to identify test data"`
}
