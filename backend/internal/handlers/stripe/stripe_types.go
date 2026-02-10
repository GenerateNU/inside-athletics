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

type CreateStripeProductRequest struct {
	Name          string        `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   string        `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
}

type CreateStripePriceRequest struct {
	product_ID    string   `json:"id" example:"product_123" doc:"ID of the product"`
	UnitAmount    float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

type UpdateStripeProductRequest struct {
	Name          *string   `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   *string   `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
}

type UpdateStripePriceRequest struct {
	UnitAmount    *float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      *Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount *int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

type GetStripeProductByIDParams struct {
	ID            string   `json:"id" example:"product_123" doc:"ID of the product"`
}

type GetStripePriceByIDParams struct {
	ID            string   `json:"id" example:"price_123" doc:"ID of the product"`
}

type GetAllStripeProductsRequest struct {
}

type GetAllStripePricesRequest struct {
	ID string `json:"id" example:"price_123" doc:"ID of the product"`
}

type GetAllStripeProductsResponse struct {
	StripeProducts []StripeProductResponse `json:"stripe_products" doc:"List of stripe products"`
	Total          int                     `json:"total" example:"25" doc:"Total number of products"`
}

type GetAllStripePricesResponse struct {
	StripePrices []StripePriceResponse `json:"stripe_products" doc:"List of stripe products"`
	Total          int                 `json:"total" example:"25" doc:"Total number of prices"`
}

type ArchiveStripeProductRequest struct {
	ID string `json:"id" example:"product_123" doc:"ID of the product"`
}

type ArchiveStripePriceRequest struct {
	ID string `json:"id" example:"price_123" doc:"ID of the product"`
}

type StripeProductResponse struct {
	ID            string   `json:"id" example:"price_123" doc:"ID of the product"`
	Name          string   `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   string   `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
}

type StripePriceResponse struct {
	ID            string   `json:"id" example:"price_123" doc:"ID of the product"`
	UnitAmount    float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

func ToStripeProductResponse(stripe_product *models.StripeProduct) *StripeProductResponse {
	return &StripeProductResponse{
		ID:            stripe_product.ID,
		Name:          stripe_product.Name,
		Description:   stripe_product.Description,
	}
}

func ToStripePriceResponse(stripe_price *models.StripePrice) *StripePriceResponse {
	return &StripePriceResponse{
		ID:            stripe_price.ID,
		UnitAmount:    stripe_price.UnitAmount,
		Interval:      Interval(stripe_price.Interval),
		IntervalCount: stripe_price.IntervalCount,
	}
}

