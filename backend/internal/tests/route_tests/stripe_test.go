package routeTests

import (
	s "inside-athletics/internal/handlers/stripe"
	"testing"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

func TestGetCustomer(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Suli"
	email := "suli@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"
	params := &stripe.CustomerParams{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}
	result, err := customer.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID

	resp := api.Get("/api/v1/stripe_customers/"+id, "Authorization: Bearer "+uuid.NewString())

	var c s.GetStripeCustomerResponse

	DecodeTo(&c, resp)

	if *c.Name != name ||
		c.Email != email ||
		*c.Phone != phone ||
		*c.Description != description {
		t.Fatalf("Unexpected response: %+v", c)
	}
}

func TestRegisterCustomer(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Suli"
	email := "suli@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"
	reqBody := s.RegisterStripeCustomerBody{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}

	resp := api.Post("/api/v1/stripe_customers/", reqBody, "Authorization: Bearer "+uuid.NewString())

	var c s.RegisterStripeCustomerResponse
	DecodeTo(&c, resp)

	if c.ID == "" {
		t.Fatalf("Expected ID to be set, got nil")
	}

	params := &stripe.CustomerParams{}
	customer, err := customer.Get(c.ID, params)

	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}
	if customer == nil {
		t.Fatalf("Unexpected")
	}

	if customer.Name != name ||
		customer.Email != email ||
		customer.Phone != phone ||
		customer.Description != description {
		t.Fatalf("Unexpected response: %+v", c)
	}
}

func TestUpdateCustomer(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Suli"
	email := "suli@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"
	params := &stripe.CustomerParams{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}
	result, err := customer.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID

	updatedName := "New Name"
	updatedEmail := "updatedemail@gmail.com"
	updatedPhone := "676 967 6967"
	updatedDecription := "updated description"
	reqBody := s.UpdateStripeCustomerBody{
		Name:        &updatedName,
		Email:       &updatedEmail,
		Phone:       &updatedPhone,
		Description: &updatedDecription,
	}

	resp := api.Patch("/api/v1/stripe_customers/"+id, reqBody, "Authorization: Bearer "+uuid.NewString())

	var c s.GetStripeCustomerResponse

	DecodeTo(&c, resp)
	if *c.Name != updatedName ||
		c.Email != updatedEmail ||
		*c.Phone != updatedPhone ||
		*c.Description != updatedDecription {
		t.Fatalf("Unexpected response: %+v", c)
	}
}
