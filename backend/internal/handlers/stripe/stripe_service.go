package stripe

import (
	"context"
	"inside-athletics/internal/utils"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/customer"
)

type StripeService struct{}

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
			ID: uuid.MustParse(result.ID),
		},
	}, nil
}

func (s *StripeService) GetStripeCustomer(ctx context.Context, input *GetStripeCustomerInput) (*utils.ResponseBody[GetStripeCustomerResponse], error) {
	id := input.ID

	params := &stripe.CustomerParams{}
	result, err := customer.Get(id.String(), params)
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

	result, err := customer.Get(id.String(), params)
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
