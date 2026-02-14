package stripe

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"fmt"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"

	"github.com/stripe/stripe-go/v72/customer"
)

type StripeService struct {
}

func NewStripeService(db *gorm.DB) *StripeService {
	return &StripeService{}
}

// products are created with default no associated prices
func (s *StripeService) CreateStripeProduct(ctx context.Context, input *struct{ Body CreateStripeProductRequest }) (*utils.ResponseBody[stripe.Product], error) {
	// Validate business rules
	if input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty.")
	}

	if input.Body.Description == "" {
		return nil, huma.Error422UnprocessableEntity("description cannot be empty.")
	}

	product_params := &stripe.ProductParams{
		Name: stripe.String(input.Body.Name),
		Description: stripe.String(input.Body.Description),
	}

	stripe_product, err := utils.HandleDBError(product.New(product_params))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[stripe.Product]{
		Body: stripe_product,
	}, nil
}

func (s *StripeService) CreateStripePrice(ctx context.Context, input *struct{ Body CreateStripePriceRequest }) (*utils.ResponseBody[stripe.Price], error) {
	// Validate business rules
	if input.Body.Product_ID == "" {
		return nil, huma.Error422UnprocessableEntity("ID cannot be empty.")
	}
	if input.Body.UnitAmount <= 0 {
		return nil, huma.Error422UnprocessableEntity("Unit amount cannot be less than or equal to 0.")
	}

	if input.Body.Interval == "" {
		return nil, huma.Error422UnprocessableEntity("Interval cannot be empty.")
	}

	if input.Body.IntervalCount <= 0 {
		return nil, huma.Error422UnprocessableEntity("Interval count cannot be less than or equal to 0.")
	}

	price_params := &stripe.PriceParams{
		Product:    stripe.String(input.Body.Product_ID),
		UnitAmount: stripe.Int64(int64(input.Body.UnitAmount) / 100), // multiply by 100 since stripe does not take floats
		Currency:   stripe.String(string(stripe.CurrencyUSD)),        //hardcoded USD
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(input.Body.Interval)),
			IntervalCount: stripe.Int64(int64(input.Body.IntervalCount)),
		},
	}

	stripe_price, err := price.New(price_params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[stripe.Price]{
		Body: stripe_price,
	}, nil
}

func (s *StripeService) GetStripeProductByID(ctx context.Context, input *GetStripeProductByIDParams) (*utils.ResponseBody[stripe.Product], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}

	stripe_product, err := product.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[stripe.Product]{
		Body: stripe_product,
	}, nil
}

func (s *StripeService) GetStripePriceByID(ctx context.Context, input *GetStripePriceByIDParams) (*utils.ResponseBody[stripe.Price], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("price ID is empty")
	}
	stripe_price, err := price.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[stripe.Price]{
		Body: stripe_price,
	}, nil
}

func (s *StripeService) UpdateStripeProduct(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body UpdateStripeProductRequest
}) (*utils.ResponseBody[stripe.Product], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}

	product_params := &stripe.ProductParams{
		Name:        input.Body.Name,
		Description: input.Body.Description,
	}

	stripe_product, err := product.Update(input.ID, product_params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[stripe.Product]{
		Body: stripe_product,
	}, nil
}

func (s *StripeService) UpdateStripePrice(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body UpdateStripePriceRequest
}) (*utils.ResponseBody[stripe.Price], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("price ID is empty")
	}

	price_params := &stripe.PriceParams{
		UnitAmount: stripe.Int64(int64(*input.Body.UnitAmount) * 100), // multiply by 100 since stripe does not take floats
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(*input.Body.Interval)),
			IntervalCount: stripe.Int64(int64(*input.Body.IntervalCount)),
		},
	}

	stripe_price, err := price.Update(input.ID, price_params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[stripe.Price]{
		Body: stripe_price,
	}, nil
}

func (s *StripeService) GetAllStripeProducts(ctx context.Context, input *GetAllStripeProductsRequest) (*utils.ResponseBody[[]*stripe.Product], error) {
	params := &stripe.ProductListParams{}
	iter := product.List(params)
	products := make([]*stripe.Product, 0)

	for iter.Next() {
		products = append(products, iter.Product())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[[]*stripe.Product]{
		Body: &products,
	}, nil
}

func (s *StripeService) GetAllStripePrices(ctx context.Context, input *GetAllStripePricesRequest) (*utils.ResponseBody[[]*stripe.Price], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}
	params := &stripe.PriceListParams{
		Product: stripe.String(input.ID),
	}

	params.Active = stripe.Bool(true)
	iter := price.List(params)
	prices := make([]*stripe.Price, 0)

	for iter.Next() {
		prices = append(prices, iter.Price())
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[[]*stripe.Price]{
		Body: &prices,
	}, nil
}

// apparently you can't delete a product due to historical billing data, but you can archive it
// archiving a product automatically archives all the prices associated with it
func (s *StripeService) ArchiveStripeProduct(ctx context.Context, input *ArchiveStripeProductRequest) (*utils.ResponseBody[*stripe.Product], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}
	params := &stripe.ProductParams{
		Active: stripe.Bool(false),
	}
	stripe_product, err := product.Update(input.ID, params)
	if err != nil {
		return nil, err
	}
	priceParams := &stripe.PriceListParams{
		Product: stripe.String(input.ID),
	}
	i := price.List(priceParams)

	for i.Next() {
		p := i.Price()
		_, err = price.Update(p.ID, &stripe.PriceParams{
			Active: stripe.Bool(false),
		})
		if err != nil {
			continue
		}
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[*stripe.Product]{
		Body: &stripe_product,
	}, nil
}

func (s *StripeService) ArchiveStripePrice(ctx context.Context, input *ArchiveStripePriceRequest) (*utils.ResponseBody[*stripe.Price], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("price ID is empty")
	}
	params := &stripe.PriceParams{
		Active: stripe.Bool(false),
	}

	stripe_price, err := price.Update(input.ID, params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[*stripe.Price]{
		Body: &stripe_price,
	}, nil
}

func (s *StripeService) RegisterStripeCustomer(ctx context.Context, input *RegisterStripeCustomerInput) (*utils.ResponseBody[RegisterStripeCustomerResponse], error) {
	params := &stripe.CustomerParams{
		Name:        input.Body.Name,
		Email:       input.Body.Email,
		Phone:       input.Body.Phone,
		Description: input.Body.Description,
	}

	result, err := customer.New(params)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[RegisterStripeCustomerResponse]{
		Body: &RegisterStripeCustomerResponse{
			ID: result.ID,
		},
	}, nil
}

func (s *StripeService) GetStripeCustomer(ctx context.Context, input *GetStripeCustomerInput) (*utils.ResponseBody[GetStripeCustomerResponse], error) {
	id := input.ID

	params := &stripe.CustomerParams{}
	result, err := customer.Get(id, params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[GetStripeCustomerResponse]{
		Body: mapStripeCustomerToModel(result),
	}, nil
}

func (s *StripeService) UpdateStripeCustomer(ctx context.Context, input *UpdateStripeCustomerInput) (*utils.ResponseBody[GetStripeCustomerResponse], error) {
	id := input.ID

	params := &stripe.CustomerParams{
		Name:          input.Body.Name,
		Email:         input.Body.Email,
		Phone:         input.Body.Phone,
		Description:   input.Body.Description,
		DefaultSource: input.Body.DefaultSource,
		Balance:       input.Body.Balance,
		Coupon:        input.Body.Coupon,
		PromotionCode: input.Body.PromotionCode,
	}

	if input.Body.TaxExempt != nil {
		params.TaxExempt = stripe.String(*input.Body.TaxExempt)
	}

	if input.Body.Metadata != nil {
		for k, v := range input.Body.Metadata {
			params.AddMetadata(k, v)
		}
	}

	result, err := customer.Update(id, params)
	if err != nil {
		return nil, err
	}
	return &utils.ResponseBody[GetStripeCustomerResponse]{
		Body: mapStripeCustomerToModel(result),
	}, nil
}

/**
func (s *StripeService) DeleteStripeCustomer(ctx context.Context, input *DeleteStripeCustomerInput) (*utils.ResponseBody[DeleteStripeCustomerResponse], error) {
	id := input.ID

	params := &stripe.CustomerParams{}

	results, err := customer.Del(id.String(), params)
	if err != nil {
		return nil, err
	}
	respBody := &DeleteStripeCustomerResponse{
		ID:      results.ID,
		Object:  results.Object,
		Deleted: results.Deleted,
	}
	return &utils.ResponseBody[DeleteStripeCustomerResponse]{
		Body: respBody,
	}, nil
}
*/

func mapStripeCustomerToModel(c *stripe.Customer) *GetStripeCustomerResponse {
	customer := &GetStripeCustomerResponse{
		ID:                  c.ID,
		Object:              string(c.Object),
		Balance:             c.Balance,
		Created:             c.Created,
		Delinquent:          c.Delinquent,
		Email:               c.Email,
		InvoicePrefix:       c.InvoicePrefix,
		Livemode:            c.Livemode,
		Metadata:            c.Metadata,
		NextInvoiceSequence: int(c.NextInvoiceSequence),
		TaxExempt:           string(c.TaxExempt),
	}

	// Map optional fields
	if c.DefaultSource != nil && c.DefaultSource.ID != "" {
		customer.DefaultSource = &c.DefaultSource.ID
	}
	if c.Description != "" {
		customer.Description = &c.Description
	}
	if c.Name != "" {
		customer.Name = &c.Name
	}
	if c.Phone != "" {
		customer.Phone = &c.Phone
	}
	if len(c.PreferredLocales) > 0 {
		customer.PreferredLocales = c.PreferredLocales
	}

	// Map invoice settings
	if c.InvoiceSettings != nil {
		customer.InvoiceSettings = &InvoiceSettings{}
		if c.InvoiceSettings.DefaultPaymentMethod != nil {
			pm := c.InvoiceSettings.DefaultPaymentMethod.ID
			customer.InvoiceSettings.DefaultPaymentMethod = &pm
		}
		if c.InvoiceSettings.Footer != "" {
			customer.InvoiceSettings.Footer = &c.InvoiceSettings.Footer
		}
	}

	return customer
}

// Checkout Sessions CRUD
