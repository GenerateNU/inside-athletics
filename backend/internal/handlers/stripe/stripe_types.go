package stripe

import (
	models "inside-athletics/internal/models"
)

type Interval string

const (
	Day   Interval = "day"
	Week  Interval = "week"
	Month Interval = "month"
	Year  Interval = "year"
)

// CreateStripeProductRequest defines the request body for creating a new product
// NOTE: stripe offers a usage_type parameter where you can bill based on usage, and it hasn't been included for now
type CreateStripeProductRequest struct {
	Name          string   `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   string   `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
	UnitAmount    float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

type UpdateStripeProductRequest struct {
	Name          *string   `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   *string   `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
	UnitAmount    *float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      *Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount *int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

type GetAllStripeProductsResponse struct {
	StripeProducts []StripeProductResponse `json:"stripe_products" doc:"List of stripe products"`
	Total          int                     `json:"total" example:"25" doc:"Total number of sports"`
}

type DeleteStripeResponseRequest struct {
	ID string `json:"id" example:"price_123" doc:"ID of the product"`
}

type StripeProductResponse struct {
	ID            string   `json:"id" example:"price_123" doc:"ID of the product"`
	Name          string   `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   string   `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
	UnitAmount    float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

// ToSportResponse converts a Sport model to a SportResponse
func ToStripeProductResponse(stripe_product *models.StripeProduct) *StripeProductResponse {
	return &StripeProductResponse{
		ID:            stripe_product.ID,
		Name:          stripe_product.Name,
		UnitAmount:    stripe_product.UnitAmount,
		Interval:      Interval(stripe_product.Interval),
		IntervalCount: stripe_product.IntervalCount,
	}
}
