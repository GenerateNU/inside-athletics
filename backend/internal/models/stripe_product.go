package models

type Interval string

const (
	Day   Interval = "day"
	Week  Interval = "week"
	Month Interval = "month"
	Year  Interval = "year"
)

type StripeProduct struct {
	ID            string   `json:"id" example:"price_123" doc:"ID of the product"`
	Name          string   `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description   string   `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
	UnitAmount    float32  `json:"total" example:"25.50" doc:"Price per billing cycle."`
	Interval      Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}