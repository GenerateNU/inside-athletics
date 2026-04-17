package utility

type AccessCheckResponse struct {
	HasPremium bool `json:"has_premium" doc:"Whether the user has premium content access"`
	IsAdmin    bool `json:"is_admin" doc:"Whether the user has admin access"`
}
