package stripe

type Interval string

const (
	Day   Interval = "day"
	Week  Interval = "week"
	Month Interval = "month"
	Year  Interval = "year"
)

type CreateStripeProductRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description string `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
}

type CreateStripePriceRequest struct {
	Product_ID    string   `json:"id" example:"product_123" doc:"ID of the product"`
	UnitAmount    int      `json:"total" example:"2550" doc:"Price per billing cycle."`
	Interval      Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

type UpdateStripeProductRequest struct {
	Name        *string `json:"name" binding:"required,min=1,max=100" example:"Premium Plan"`
	Description *string `json:"description" binding:"required,min=1,max=200" example:"Get premium content with this subscription"`
}

type CreateStripeCheckoutSessionRequest struct {
	PriceID    string `json:"price_id" binding:"required" example:"price_12345"`
	SuccessURL string `json:"success_url" binding:"required,url" example:"https://example.com/success"`
	CancelURL  string `json:"cancel_url" binding:"required,url" example:"https://example.com/cancel"`
	Quantity   int64  `json:"quantity" binding:"required,min=1" example:"1"`
}

type UpdateStripePriceRequest struct {
	UnitAmount    *int      `json:"total" example:"2550" doc:"Price per billing cycle."`
	Interval      *Interval `json:"interval" example:"day" doc:"Interval between payments"`
	IntervalCount *int      `json:"interval_count" example:"3" doc:"Number of intervals a billing cycle lasts"`
}

type GetStripeProductByIDParams struct {
	ID string `path:"id" example:"prod_123" doc:"ID of the product"`
}

type GetStripePriceByIDParams struct {
	ID string `path:"id" example:"prod_123" doc:"ID of the product"`
}

type GetStripeCheckoutSessionParams struct {
	ID string `path:"id" example:"prod_123" doc:"ID of the product"`
}

type GetAllStripeProductsRequest struct {
}

type GetAllStripePricesRequest struct {
	ID string `json:"id" example:"price_123" doc:"ID of the product"`
}

type GetAllStripeSessionsRequest struct {
	CustomerID string `query:"customer_id" doc:"Filter by customer"`
	Limit      int64  `query:"limit" doc:"Number of sessions to return"`
}

type ArchiveStripeProductRequest struct {
	ID string `json:"id" example:"product_123" doc:"ID of the product"`
}

type ArchiveStripePriceRequest struct {
	ID string `json:"id" example:"price_123" doc:"ID of the product"`
}

type DeleteCheckoutSessionRequest struct {
	ID string `path:"id" example:"cs_123" doc:"ID of the checkout session"`
}

type GetStripeCustomerInput struct {
	ID string `path:"id" maxLength:"36" example:"1" doc:"ID to identify the stripe user"`
}

type GetStripeCustomerResponse struct {
	ID                  string            `gorm:"primaryKey;column:id" json:"id" example:"cus_NffrFeUfNV2Hib" doc:"Stripe customer ID"`
	Object              string            `gorm:"column:object" json:"object" example:"customer" doc:"Object type"`
	Address             *string           `gorm:"column:address;type:jsonb" json:"address" doc:"Customer address"`
	Balance             int64             `gorm:"column:balance;default:0" json:"balance" example:"0" doc:"Customer balance in cents"`
	Created             int64             `gorm:"column:created" json:"created" example:"1680893993" doc:"Unix timestamp of creation"`
	Currency            *string           `gorm:"column:currency" json:"currency" doc:"Default currency"`
	DefaultSource       *string           `gorm:"column:default_source" json:"default_source" doc:"Default payment source ID"`
	Delinquent          bool              `gorm:"column:delinquent;default:false" json:"delinquent" example:"false" doc:"Whether customer is delinquent"`
	Description         *string           `gorm:"column:description" json:"description" doc:"Customer description"`
	Email               string            `gorm:"column:email" json:"email" example:"jennyrosen@example.com" doc:"Customer email"`
	InvoicePrefix       string            `gorm:"column:invoice_prefix" json:"invoice_prefix" example:"0759376C" doc:"Invoice number prefix"`
	InvoiceSettings     *InvoiceSettings  `gorm:"column:invoice_settings;type:jsonb;serializer:json" json:"invoice_settings" doc:"Invoice settings"`
	Livemode            bool              `gorm:"column:livemode;default:false" json:"livemode" example:"false" doc:"Whether in live mode"`
	Metadata            map[string]string `gorm:"column:metadata;type:jsonb;serializer:json" json:"metadata" doc:"Custom metadata"`
	Name                *string           `gorm:"column:name" json:"name" example:"Jenny Rosen" doc:"Customer name"`
	NextInvoiceSequence int               `gorm:"column:next_invoice_sequence;default:1" json:"next_invoice_sequence" example:"1" doc:"Next invoice sequence number"`
	Phone               *string           `gorm:"column:phone" json:"phone" doc:"Customer phone number"`
	PreferredLocales    []string          `gorm:"column:preferred_locales;type:jsonb;serializer:json" json:"preferred_locales" doc:"Preferred locales"`
	Shipping            *string           `gorm:"column:shipping;type:jsonb" json:"shipping" doc:"Shipping information"`
	TaxExempt           string            `gorm:"column:tax_exempt;default:'none'" json:"tax_exempt" example:"none" doc:"Tax exempt status"`
	TestClock           *string           `gorm:"column:test_clock" json:"test_clock" doc:"Test clock ID"`
}

type GetStripeCustomerByEmailInput struct {
	Email string `path:"email" maxLength:"36" example:"suli@gmail.com" doc:"email to identify the stripe user"`
}

type GetStripeCustomerByEmailResponse struct {
	ID    string `gorm:"primaryKey;column:id" json:"id" example:"cus_NffrFeUfNV2Hib" doc:"Stripe customer ID"`
	Email string `gorm:"column:email" json:"email" example:"jennyrosen@example.com" doc:"Customer email"`
}

type InvoiceSettings struct {
	CustomFields         *string `json:"custom_fields"`
	DefaultPaymentMethod *string `json:"default_payment_method"`
	Footer               *string `json:"footer"`
	RenderingOptions     *string `json:"rendering_options"`
}

type RegisterStripeCustomerInput struct {
	Body RegisterStripeCustomerBody
}

type RegisterStripeCustomerBody struct {
	Name        *string `json:"name" example:"Suli" doc:"The name of the user"`
	Email       *string `json:"email" example:"suli@northeastern.edu" doc:"The email of the user"`
	Phone       *string `json:"phone" example:"(111) 222-3333" doc:"The phone number of the user"`
	Description *string `json:"description" example:"A verified athelete" doc:"A description of the user"`
}

type RegisterStripeCustomerResponse struct {
	ID string `json:"id" example:"1" doc:"ID of the user on stripe"`
}

type UpdateStripeCustomerInput struct {
	ID   string `path:"id" maxLength:"36" example:"1" doc:"ID to identify the stripe user"`
	Body UpdateStripeCustomerBody
}

type UpdateStripeCustomerBody struct {
	Balance          *int64            `json:"balance,omitempty" doc:"Current balance in cents (can be negative)"`
	Coupon           *string           `json:"coupon,omitempty" doc:"Coupon ID to apply discount"`
	DefaultSource    *string           `json:"default_source,omitempty" doc:"ID of default payment source"`
	Description      *string           `json:"description,omitempty" doc:"Description of customer"`
	Email            *string           `json:"email,omitempty" example:"newemail@example.com" doc:"Customer email"`
	Metadata         map[string]string `json:"metadata,omitempty" doc:"Set of key-value pairs for metadata"`
	Name             *string           `json:"name,omitempty" example:"Jane Doe" doc:"Customer name"`
	Phone            *string           `json:"phone,omitempty" example:"+15555551234" doc:"Customer phone number"`
	PreferredLocales []string          `json:"preferred_locales,omitempty" doc:"Customer's preferred languages"`
	PromotionCode    *string           `json:"promotion_code,omitempty" doc:"Promotion code to apply"`
	TaxExempt        *string           `json:"tax_exempt,omitempty" example:"none" doc:"Tax exemption status: none, exempt, or reverse"`
}

type CustomField struct {
	Name  string `json:"name" example:"Tax ID" doc:"Field name"`
	Value string `json:"value" example:"123-45-6789" doc:"Field value"`
}

type DeleteStripeCustomerInput struct {
	ID string `json:"path" maxLength:"36" example:"1" doc:"ID of the user on stripe"`
}

type DeleteStripeCustomerResponse struct {
	ID      string `gorm:"primaryKey;column:id" json:"id" example:"cus_NffrFeUfNV2Hib" doc:"Stripe customer ID"`
	Object  string `gorm:"column:object" json:"object" example:"customer" doc:"Object type"`
	Deleted bool   `gorm:"column:deleted" json:"deleted" example:"true" doc:"Object bool"`
}

type HasActiveSubscriptionInput struct {
	CustomerID string `json:"path" doc:"Stripe Customer ID"`
}

type HasActiveSubscriptionResponse struct {
	HasActiveSubscription bool   `json:"has_active_subscription"`
	SubscriptionID        string `json:"subscription_id,omitempty"`
	Status                string `json:"status,omitempty"` // active, past_due, canceled, etc
	CurrentPeriodEnd      int64  `json:"current_period_end,omitempty"`
}
