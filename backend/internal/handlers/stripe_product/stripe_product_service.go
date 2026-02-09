package stripe_product

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/price"
)

type StripeProductService struct {
}

// NewStripeProductService creates a new SportService instance
func NewStripeProductService(db *gorm.DB) *StripeProductService {
	return &StripeProductService{
		stripeServiceDB: NewStripeServiceDB(db),
	}
}

func (s *StripeProductService) CreateStripeProduct(ctx context.Context, input *struct{ Body CreateStripeProductRequest }) (*utils.ResponseBody[StripeProductResponse], error) {
	// Validate business rules
	if input.Body.Name == "" {
		return nil, huma.Error422UnprocessableEntity("name cannot be empty.")
	}

	if input.Body.Description == "" {
		return nil, huma.Error422UnprocessableEntity("description cannot be empty.")
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

	product_params := &stripe.ProductParams{
    Name: stripe.String(input.Body.Name),
	}

	stripe_product, err := utils.HandleDBError(product.New(product_params))
	if err != nil {
		return nil, err
	}

	price_params = &stripe.PriceParams{
		Product:    stripe.String(stripe_product.ID),
		UnitAmount: stripe.Int64(int64(input.Body.UnitAmount) * 100), // multiply by 100 since stripe does not take floats
		Currency:   stripe.String(string(stripe.CurrencyUSD)), //hardcoded USD
		Recurring: &stripe.PriceRecurringParams {
			Interval: stripe.String(string(input.Body.Interval)),
			IntervalCount: stripe.Int64(int64(input.Body.IntervalCount))
		}
	}

	stripe_price, err := price.New(price_params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(stripe_product),
	}, nil
}


func (s *StripeProductService) GetStripeProductByID(ctx context.Context, input *GetStripeProductByIDParams) (*utils.ResponseBody[StripeProductResponse], error) {
	stripe_product, err := utils.HandleDBError(product.Get(input.Body.ID, nil))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(stripe_product),
	}, nil
}

// UpdatePost updates an existing post
func (s *StripeProductService) UpdateStripeProduct(id String, updates UpdateStripeProductRequest) (*models.StripeProduct, error) {
	params := &stripe.ProductParams{
    	Name: stripe.String(input.Body.Name),
		Description: stripe.String(input.Body.Description),
		UnitAmount: stripe.int(input.Body.UnitAmount),
		Recurring: &stripe.PriceRecurringParams {
			Interval: stripe.String(input.Body.Interval),
			IntervalCount: stripe.String(input.Body.IntervalCount),
		}
	}
	
	stripe_product, err := product.Update("prod_NWjs8kKbJWmuuc", params)

	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(stripe_product),
	}, nil
}

func (s *StripeProductService) GetAllStripeProducts(ctx context.Context, input *GetAllStripeProductParams) (*utils.ResponseBody[GetAllStripeProductsResponse], error) {
	stripe_products := stripe.Product.list()
	stripeProductResponses := make([]stripeProductsResponse, 0, len(sports))

	for iter.Next() {
		prod := iter.Product()
		stripeProductResponses = append(prod, *ToStripeProductResponse(&prod))
	}

	if err := iter.Err(); err != nil {
    log.Fatalf("Error during pagination: %v", err)
	}

	return &utils.ResponseBody[GetAllStripeProductsResponse]{
		Body: &GetAllStripeProductsResponse{
			StripeProducts: stripeProductResponses,
			Total:  int(total),
		},
	}, nil
}

func (s *StripeProductService) DeleteStripeProduct(ctx context.Context, input *DeleteStripeProductRequest) (*utils.ResponseBody[StripeProductResponse], error) {
	stripe_product, err := utils.HandleDBError(stripe.product.Del(input.ID, nil))
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: ToStripeProductResponse(stripe_product),
	}, nil
}
