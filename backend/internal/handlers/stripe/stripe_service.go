package stripe

import (
	"context"
	models "inside-athletics/internal/models"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"

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
func (s *StripeService) CreateStripeProduct(ctx context.Context, input *struct{ Body CreateStripeProductRequest }) (*utils.ResponseBody[StripeProductResponse], error) {
	// Validate business rules
	if input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty.")
	}

	if input.Body.Description == "" {
		return nil, huma.Error422UnprocessableEntity("description cannot be empty.")
	}

	product_params := &stripe.ProductParams{
		Name: stripe.String(input.Body.Name),
	}

	stripe_product, err := utils.HandleDBError(product.New(product_params))
	if err != nil {
		return nil, err
	}

	final := models.StripeProduct{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(&final),
	}, nil
}

func (s *StripeService) CreateStripePrice(ctx context.Context, input *struct{ Body CreateStripePriceRequest }) (*utils.ResponseBody[StripePriceResponse], error) {
	// Validate business rules
	if input.Body.product_ID == "" {
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
		Product:    stripe.String(input.Body.product_ID),
		UnitAmount: stripe.Int64(int64(input.Body.UnitAmount) * 100), // multiply by 100 since stripe does not take floats
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

	final := models.StripePrice{
		ID:            stripe_price.Product.ID,
		UnitAmount:    float32(stripe_price.UnitAmount),
		Interval:      models.Interval(stripe_price.Recurring.Interval),
		IntervalCount: int(stripe_price.Recurring.IntervalCount),
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: ToStripePriceResponse(&final),
	}, nil
}

func (s *StripeService) GetStripeProductByID(ctx context.Context, input *GetStripeProductByIDParams) (*utils.ResponseBody[StripeProductResponse], error) {
	stripe_product, err := product.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	priceParams := &stripe.PriceListParams{
		Product: stripe.String(stripe_product.ID),
	}
	priceParams.Filters.AddFilter("limit", "", "100")
	priceIter := price.List(priceParams)

	var prices []models.StripePrice
	for priceIter.Next() {
		p := priceIter.Price()
		prices = append(prices, models.StripePrice{
			ID:            p.ID,
			UnitAmount:    float32(p.UnitAmount) / 100.0, // Convert cents to float
			Currency:      string(p.Currency),
			Interval:      models.Interval(p.Recurring.Interval),
			IntervalCount: int(p.Recurring.IntervalCount),
		})
	}

	final := models.StripeProduct{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
		Prices:      prices,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(&final),
	}, nil
}

func (s *StripeService) GetStripePriceByID(ctx context.Context, input *GetStripePriceByIDParams) (*utils.ResponseBody[StripePriceResponse], error) {
	stripe_price, err := price.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	final := models.StripePrice{
		ID:            stripe_price.ID,
		UnitAmount:    float32(stripe_price.UnitAmount),
		Interval:      models.Interval(stripe_price.Recurring.Interval),
		IntervalCount: int(stripe_price.Recurring.IntervalCount),
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: ToStripePriceResponse(&final),
	}, nil
}

func (s *StripeService) UpdateStripeProduct(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body UpdateStripeProductRequest
}) (*utils.ResponseBody[StripeProductResponse], error) {

	product_params := &stripe.ProductParams{
		Name:        input.Body.Name,
		Description: input.Body.Description,
	}

	stripe_product, err := product.Update(input.ID, product_params)
	if err != nil {
		return nil, err
	}

	final := models.StripeProduct{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(&final),
	}, nil
}

func (s *StripeService) UpdateStripePrice(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body UpdateStripePriceRequest
}) (*utils.ResponseBody[StripePriceResponse], error) {

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

	final := models.StripePrice{
		ID:            stripe_price.ID,
		UnitAmount:    float32(stripe_price.UnitAmount),
		Interval:      models.Interval(stripe_price.Recurring.Interval),
		IntervalCount: int(stripe_price.Recurring.IntervalCount),
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: ToStripePriceResponse(&final),
	}, nil
}

func (s *StripeService) GetAllStripeProducts(ctx context.Context, input *GetAllStripeProductsRequest) (*utils.ResponseBody[GetAllStripeProductsResponse], error) {
	params := &stripe.ProductListParams{}
	params.AddExpand("data.default_price")

	iter := product.List(params)
	stripeProductResponses := make([]StripeProductResponse, 0)

	for iter.Next() {
		prod := iter.Product()
		var prices []models.StripePrice
		if prod.DefaultPrice != nil {
			p := prod.DefaultPrice
			priceItem := models.StripePrice{
				ID:         p.ID,
				UnitAmount: float32(p.UnitAmount) / 100.0,
				Currency:   string(p.Currency),
			}
			if p.Recurring != nil {
				priceItem.Interval = models.Interval(p.Recurring.Interval)
				priceItem.IntervalCount = int(p.Recurring.IntervalCount)
			}
			prices = append(prices, priceItem)
		}
		stripeProductResponses = append(stripeProductResponses, *ToStripeProductResponse(&models.StripeProduct{
			ID:          prod.ID,
			Name:        prod.Name,
			Description: prod.Description,
			Prices:      prices,
		}))
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[GetAllStripeProductsResponse]{
		Body: &GetAllStripeProductsResponse{
			StripeProducts: stripeProductResponses,
			Total:          len(stripeProductResponses),
		},
	}, nil
}

func (s *StripeService) GetAllStripePrices(ctx context.Context, input *GetAllStripePricesRequest) (*utils.ResponseBody[GetAllStripePricesResponse], error) {
	params := &stripe.PriceListParams{
		Product: stripe.String(input.ID),
	}

	params.Active = stripe.Bool(true)
	iter := price.List(params)
	stripePriceResponses := make([]StripePriceResponse, 0)

	for iter.Next() {
		p := iter.Price()
		priceModel := models.StripePrice{
			ID:         p.ID,
			UnitAmount: float32(p.UnitAmount) / 100.0,
			Currency:   string(p.Currency),
		}
		if p.Recurring != nil {
			priceModel.Interval = models.Interval(p.Recurring.Interval)
			priceModel.IntervalCount = int(p.Recurring.IntervalCount)
		}
		stripePriceResponses = append(stripePriceResponses, *ToStripePriceResponse(&priceModel))
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[GetAllStripePricesResponse]{
		Body: &GetAllStripePricesResponse{
			StripePrices: stripePriceResponses,
			Total:        len(stripePriceResponses),
		},
	}, nil
}

// apparently you can't delete a product due to historical billing data, but you can archive it
// archiving a product automatically archives all the prices associated with it
func (s *StripeService) ArchiveStripeProduct(ctx context.Context, input *ArchiveStripeProductRequest) (*utils.ResponseBody[StripeProductResponse], error) {
	params := &stripe.ProductParams{
		Active: stripe.Bool(false),
	}
	stripeProduct, err := product.Update(input.ID, params)
	if err != nil {
		return nil, err
	}
	priceParams := &stripe.PriceListParams{
		Product: stripe.String(input.ID),
	}
	i := price.List(priceParams)

	var archivedPrices []models.StripePrice

	for i.Next() {
		p := i.Price()
		updatedPrice, err := price.Update(p.ID, &stripe.PriceParams{
			Active: stripe.Bool(false),
		})
		if err != nil {
			continue
		}
		priceModel := models.StripePrice{
			ID:         updatedPrice.ID,
			UnitAmount: float32(updatedPrice.UnitAmount) / 100.0,
			Currency:   string(updatedPrice.Currency),
		}
		if updatedPrice.Recurring != nil {
			priceModel.Interval = models.Interval(updatedPrice.Recurring.Interval)
			priceModel.IntervalCount = int(updatedPrice.Recurring.IntervalCount)
		}

		archivedPrices = append(archivedPrices, priceModel)
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	final := models.StripeProduct{
		ID:          stripeProduct.ID,
		Name:        stripeProduct.Name,
		Description: stripeProduct.Description,
		Prices:      archivedPrices,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(&final),
	}, nil
}

func (s *StripeService) ArchiveStripePrice(ctx context.Context, input *ArchiveStripePriceRequest) (*utils.ResponseBody[StripePriceResponse], error) {
	params := &stripe.PriceParams{
		Active: stripe.Bool(false),
	}

	updatedPrice, err := price.Update(input.ID, params)
	if err != nil {
		return nil, err
	}

	final := models.StripePrice{
		ID:         updatedPrice.ID,
		UnitAmount: float32(updatedPrice.UnitAmount) / 100.0,
		Currency:   string(updatedPrice.Currency),
	}

	if updatedPrice.Recurring != nil {
		final.Interval = models.Interval(updatedPrice.Recurring.Interval)
		final.IntervalCount = int(updatedPrice.Recurring.IntervalCount)
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: ToStripePriceResponse(&final),
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

func (s *StripeService) DeleteStripeCustomer(ctx context.Context, input *DeleteStripeCustomerInput) (*utils.ResponseBody[DeleteStripeCustomerResponse], error) {
	id := input.ID

	params := &stripe.CustomerParams{}

	results, err := customer.Del(id, params)
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
