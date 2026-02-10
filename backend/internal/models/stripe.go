package models

type Interval string

const (
	Day   Interval = "day"
	Week  Interval = "week"
	Month Interval = "month"
	Year  Interval = "year"
)

type StripePrice struct {
	ID            string   `json:"id" example:"price_123" doc:"ID of the price"`
	UnitAmount    float32  `json:"unit_amount" example:"25.50" doc:"Price per billing cycle."`
	Currency      string   `json:"currency" example:"usd" doc:"Three-letter ISO currency code"`
	Interval      Interval `json:"interval" example:"month" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"1" doc:"Number of intervals a billing cycle lasts"`
}

type StripeProduct struct {
	ID          string        `json:"id" example:"prod_123" doc:"ID of the product"`
	Name        string        `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description string        `json:"description" binding:"required,min=1,max=200" example:"Get premium content"`
	Prices      []StripePrice `json:"prices" doc:"List of prices associated with this product"`
}