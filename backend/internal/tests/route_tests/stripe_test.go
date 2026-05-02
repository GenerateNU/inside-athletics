package routeTests

import (
	s "inside-athletics/internal/handlers/stripe"
	"testing"

	"github.com/google/uuid"
)

func TestCreateProduct(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	resp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name:        "Premium Plan",
		Description: "Get premium content with this subscription",
	}, "Authorization: Bearer "+uuid.NewString())
	var prod s.StripeProductResponse
	DecodeTo(&prod, resp)

	if prod.Name != "Premium Plan" {
		t.Errorf("expected Name to be 'Premium Plan', got %s", prod.Name)
	}
	if prod.Description != "Get premium content with this subscription" {
		t.Errorf("expected Description to be 'Get premium content...', got %s", prod.Description)
	}
}

func TestGetProductByID(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	name := "Premium Plan"
	description := "Get premium content with this subscription"

	createResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name:        name,
		Description: description,
	}, auth)
	var created s.StripeProductResponse
	DecodeTo(&created, createResp)

	resp := api.Get(
		"/api/v1/stripe_product/"+created.ID,
		s.GetStripeProductByIDParams{ID: created.ID},
		auth,
	)
	var prod s.StripeProductResponse
	DecodeTo(&prod, resp)

	if prod.Name != name {
		t.Errorf("expected Name to be '%s', got %s", name, prod.Name)
	}
	if prod.Description != description {
		t.Errorf("expected Description to be '%s', got %s", description, prod.Description)
	}
}

func TestUpdateStripeProduct(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name:        "Premium Plan",
		Description: "Get premium content with this subscription",
	}, auth)
	var created s.StripeProductResponse
	DecodeTo(&created, createResp)

	newName := "Basic Plan"
	newDescription := "The plan you get when you wanna miss out on premium content..."
	resp := api.Patch(
		"/api/v1/stripe_product/"+created.ID,
		s.UpdateStripeProductRequest{Name: &newName, Description: &newDescription},
		auth,
	)
	var prod s.StripeProductResponse
	DecodeTo(&prod, resp)

	if prod.Name != newName {
		t.Errorf("expected Name to be '%s', got %s", newName, prod.Name)
	}
	if prod.Description != newDescription {
		t.Errorf("expected Description to be '%s', got %s", newDescription, prod.Description)
	}
}

func TestArchiveStripeProduct(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name:        "Premium Plan",
		Description: "Get premium content with this subscription",
	}, auth)
	var created s.StripeProductResponse
	DecodeTo(&created, createResp)

	resp := api.Delete(
		"/api/v1/stripe_product/"+created.ID,
		s.ArchiveStripeProductRequest{ID: created.ID},
		auth,
	)
	var prod s.StripeProductResponse
	DecodeTo(&prod, resp)

	if prod.Active {
		t.Errorf("expected product to be inactive after delete, but Active = true")
	}
}

func TestGetAllProducts(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Basic Plan", Description: "Get basic content with this subscription",
	}, auth)

	resp := api.Get("/api/v1/stripe_products/", s.GetAllStripeProductsRequest{}, auth)
	var products []s.StripeProductResponse
	DecodeTo(&products, resp)

	if len(products) < 2 {
		t.Errorf("expected at least 2 products, got %d", len(products))
	}
}

func TestCreatePrice(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	unitAmount := 2500
	interval := s.Day
	intervalCount := 3

	resp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    unitAmount,
		Interval:      interval,
		IntervalCount: intervalCount,
	}, auth)
	var pr s.StripePriceResponse
	DecodeTo(&pr, resp)

	if pr.ID == "" {
		t.Errorf("expected price ID to be set, got empty string")
	}
	if pr.ProductID != prod.ID {
		t.Errorf("expected product id to be %s, got %s", prod.ID, pr.ProductID)
	}
	if pr.UnitAmount != int64(unitAmount) {
		t.Errorf("expected unit amount to be %d, got %d", unitAmount, pr.UnitAmount)
	}
	if pr.Interval != string(interval) {
		t.Errorf("expected interval to be %s, got %s", interval, pr.Interval)
	}
	if pr.IntervalCount != int64(intervalCount) {
		t.Errorf("expected interval count to be %d, got %d", intervalCount, pr.IntervalCount)
	}
}

func TestGetStripePriceByID(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	unitAmount := 2550
	interval := s.Day
	intervalCount := 3

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    unitAmount,
		Interval:      interval,
		IntervalCount: intervalCount,
	}, auth)
	var createdPrice s.StripePriceResponse
	DecodeTo(&createdPrice, createPriceResp)

	resp := api.Get("/api/v1/stripe_price/"+createdPrice.ID, s.GetStripePriceByIDParams{ID: createdPrice.ID}, auth)
	var pr s.StripePriceResponse
	DecodeTo(&pr, resp)

	if pr.ID == "" {
		t.Errorf("expected price ID to be set, got empty string")
	}
	if pr.ProductID != prod.ID {
		t.Errorf("expected product id to be %s, got %s", prod.ID, pr.ProductID)
	}
	if pr.UnitAmount != int64(unitAmount) {
		t.Errorf("expected unit amount to be %d, got %d", unitAmount, pr.UnitAmount)
	}
	if pr.Interval != string(interval) {
		t.Errorf("expected interval to be %s, got %s", interval, pr.Interval)
	}
	if pr.IntervalCount != int64(intervalCount) {
		t.Errorf("expected interval count to be %d, got %d", intervalCount, pr.IntervalCount)
	}
}

func TestUpdateStripePrice(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    2550,
		Interval:      s.Day,
		IntervalCount: 3,
	}, auth)
	var originalPrice s.StripePriceResponse
	DecodeTo(&originalPrice, createPriceResp)

	newUnitAmount := 3000
	newInterval := s.Week
	newIntervalCount := 2

	resp := api.Patch("/api/v1/stripe_price/"+originalPrice.ID, s.UpdateStripePriceRequest{
		UnitAmount:    &newUnitAmount,
		Interval:      &newInterval,
		IntervalCount: &newIntervalCount,
	}, auth)
	var updatedPrice s.StripePriceResponse
	DecodeTo(&updatedPrice, resp)

	if updatedPrice.ID == originalPrice.ID {
		t.Errorf("expected new price ID, but got same ID %s", updatedPrice.ID)
	}
	if updatedPrice.ProductID != prod.ID {
		t.Errorf("expected product id to be %s, got %s", prod.ID, updatedPrice.ProductID)
	}
	if updatedPrice.UnitAmount != int64(newUnitAmount) {
		t.Errorf("expected unit amount to be %d, got %d", newUnitAmount, updatedPrice.UnitAmount)
	}
	if updatedPrice.Interval != string(newInterval) {
		t.Errorf("expected interval to be %s, got %s", newInterval, updatedPrice.Interval)
	}
	if updatedPrice.IntervalCount != int64(newIntervalCount) {
		t.Errorf("expected interval count to be %d, got %d", newIntervalCount, updatedPrice.IntervalCount)
	}
}

func TestArchiveStripePrice(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    2550,
		Interval:      s.Day,
		IntervalCount: 3,
	}, auth)
	var originalPrice s.StripePriceResponse
	DecodeTo(&originalPrice, createPriceResp)

	resp := api.Delete(
		"/api/v1/stripe_price/"+originalPrice.ID,
		s.ArchiveStripePriceRequest{ID: originalPrice.ID},
		auth,
	)
	var archivedPrice s.StripePriceResponse
	DecodeTo(&archivedPrice, resp)

	if archivedPrice.Active {
		t.Errorf("expected price to be inactive after delete, but Active = true")
	}
}

func TestGetAllStripePrices(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	priceBody := s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    2550,
		Interval:      s.Day,
		IntervalCount: 3,
	}
	api.Post("/api/v1/stripe_price/", priceBody, auth)
	api.Post("/api/v1/stripe_price/", priceBody, auth)

	resp := api.Get("/api/v1/stripe_prices/"+prod.ID, s.GetAllStripePricesRequest{ID: prod.ID}, auth)
	var prices []s.StripePriceResponse
	DecodeTo(&prices, resp)

	if len(prices) < 2 {
		t.Errorf("expected at least 2 prices, got %d", len(prices))
	}
}

func TestGetCustomer(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	name := "Suli"
	email := "suli@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"

	createResp := api.Post("/api/v1/stripe_customers/", s.RegisterStripeCustomerBody{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}, auth)
	var created s.RegisterStripeCustomerResponse
	DecodeTo(&created, createResp)

	resp := api.Get("/api/v1/stripe_customers/"+created.ID, auth)
	var c s.GetStripeCustomerResponse
	DecodeTo(&c, resp)

	if c.Name == nil || *c.Name != name ||
		c.Email != email ||
		c.Phone == nil || *c.Phone != phone ||
		c.Description == nil || *c.Description != description {
		t.Fatalf("Unexpected response: %+v", c)
	}
}

func TestGetCustomerByEmail(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	name := "Suli"
	email := "suli_newemail@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"

	api.Post("/api/v1/stripe_customers/", s.RegisterStripeCustomerBody{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}, auth)

	resp := api.Get("/api/v1/stripe_customers/email/"+email, auth)
	var c s.GetStripeCustomerByEmailResponse
	DecodeTo(&c, resp)

	if c.Email != email {
		t.Fatalf("Unexpected response, email does not match: %+v", c)
	}
}

func TestRegisterCustomer(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	name := "Suli"
	email := "suli@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"

	resp := api.Post("/api/v1/stripe_customers/", s.RegisterStripeCustomerBody{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}, auth)
	var c s.RegisterStripeCustomerResponse
	DecodeTo(&c, resp)

	if c.ID == "" {
		t.Fatalf("Expected ID to be set, got nil")
	}

	getResp := api.Get("/api/v1/stripe_customers/"+c.ID, auth)
	var customer s.GetStripeCustomerResponse
	DecodeTo(&customer, getResp)

	if customer.Name == nil || *customer.Name != name ||
		customer.Email != email ||
		customer.Phone == nil || *customer.Phone != phone ||
		customer.Description == nil || *customer.Description != description {
		t.Fatalf("Unexpected response: %+v", customer)
	}
}

func TestUpdateCustomer(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	auth := "Authorization: Bearer " + uuid.NewString()
	name := "Suli"
	email := "suli@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"

	createResp := api.Post("/api/v1/stripe_customers/", s.RegisterStripeCustomerBody{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}, auth)
	var created s.RegisterStripeCustomerResponse
	DecodeTo(&created, createResp)

	updatedName := "New Name"
	updatedEmail := "updatedemail@gmail.com"
	updatedPhone := "676 967 6967"
	updatedDescription := "updated description"

	resp := api.Patch("/api/v1/stripe_customers/"+created.ID, s.UpdateStripeCustomerBody{
		Name:        &updatedName,
		Email:       &updatedEmail,
		Phone:       &updatedPhone,
		Description: &updatedDescription,
	}, auth)
	var c s.GetStripeCustomerResponse
	DecodeTo(&c, resp)

	if c.Name == nil || *c.Name != updatedName ||
		c.Email != updatedEmail ||
		c.Phone == nil || *c.Phone != updatedPhone ||
		c.Description == nil || *c.Description != updatedDescription {
		t.Fatalf("Unexpected response: %+v", c)
	}
}

func TestCreateCheckoutSession(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user, _ := seedUserAndCollege(t, testDB, "create-checkout-session")
	auth := "Authorization: Bearer " + uuid.NewString()

	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    2550,
		Interval:      s.Day,
		IntervalCount: 3,
	}, auth)
	var pr s.StripePriceResponse
	DecodeTo(&pr, createPriceResp)

	resp := api.Post("/api/v1/checkout/sessions/", s.CreateStripeCheckoutSessionRequest{
		UserID:     user.ID.String(),
		PriceID:    pr.ID,
		SuccessURL: "https://example.com/success",
		CancelURL:  "https://example.com/cancel",
		Quantity:   2,
	}, auth)
	var session s.StripeCheckoutSessionResponse
	DecodeTo(&session, resp)

	if session.ID == "" {
		t.Errorf("expected checkout session ID to be set, got empty string")
	}
	if session.URL == "" {
		t.Errorf("expected checkout session URL to be set, got empty string")
	}
	if session.Mode != "subscription" {
		t.Errorf("expected mode to be subscription, got %s", session.Mode)
	}
}

func TestGetStripeCheckoutSessionByID(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user, _ := seedUserAndCollege(t, testDB, "get-checkout-session-by-id")
	auth := "Authorization: Bearer " + uuid.NewString()

	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    2550,
		Interval:      s.Day,
		IntervalCount: 3,
	}, auth)
	var pr s.StripePriceResponse
	DecodeTo(&pr, createPriceResp)

	createSessionResp := api.Post("/api/v1/checkout/sessions/", s.CreateStripeCheckoutSessionRequest{
		UserID:     user.ID.String(),
		PriceID:    pr.ID,
		SuccessURL: "https://example.com/success",
		CancelURL:  "https://example.com/cancel",
		Quantity:   1,
	}, auth)
	var createdSession s.StripeCheckoutSessionResponse
	DecodeTo(&createdSession, createSessionResp)

	resp := api.Get(
		"/api/v1/checkout/sessions/"+createdSession.ID,
		s.GetStripeCheckoutSessionParams{ID: createdSession.ID},
		auth,
	)
	var session s.StripeCheckoutSessionResponse
	DecodeTo(&session, resp)

	if session.ID == "" {
		t.Errorf("expected checkout session ID to be set, got empty string")
	}
	if session.URL == "" {
		t.Errorf("expected checkout session URL to be set, got empty string")
	}
	if session.Mode != "subscription" {
		t.Errorf("expected mode to be subscription, got %s", session.Mode)
	}
}

func TestDeleteStripeCheckoutSession(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user, _ := seedUserAndCollege(t, testDB, "delete-checkout-session")
	auth := "Authorization: Bearer " + uuid.NewString()

	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Premium Plan", Description: "Get premium content with this subscription",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    1000,
		Interval:      s.Month,
		IntervalCount: 1,
	}, auth)
	var pr s.StripePriceResponse
	DecodeTo(&pr, createPriceResp)

	createSessionResp := api.Post("/api/v1/checkout/sessions/", s.CreateStripeCheckoutSessionRequest{
		UserID:     user.ID.String(),
		PriceID:    pr.ID,
		SuccessURL: "https://example.com/success",
		CancelURL:  "https://example.com/cancel",
		Quantity:   1,
	}, auth)
	var createdSession s.StripeCheckoutSessionResponse
	DecodeTo(&createdSession, createSessionResp)

	resp := api.Delete(
		"/api/v1/checkout/sessions/"+createdSession.ID,
		s.DeleteCheckoutSessionRequest{ID: createdSession.ID},
		auth,
	)
	var session s.StripeCheckoutSessionResponse
	DecodeTo(&session, resp)

	if session.ID == "" {
		t.Errorf("expected checkout session ID to be set, got empty string")
	}
}

func TestGetAllStripeSessions(t *testing.T) {
	t.Parallel()
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	user, _ := seedUserAndCollege(t, testDB, "get-all-sessions")
	auth := "Authorization: Bearer " + uuid.NewString()

	createProdResp := api.Post("/api/v1/stripe_product/", s.CreateStripeProductRequest{
		Name: "Standard Plan", Description: "Standard subscription plan",
	}, auth)
	var prod s.StripeProductResponse
	DecodeTo(&prod, createProdResp)

	createPriceResp := api.Post("/api/v1/stripe_price/", s.CreateStripePriceRequest{
		Product_ID:    prod.ID,
		UnitAmount:    1000,
		Interval:      s.Month,
		IntervalCount: 1,
	}, auth)
	var pr s.StripePriceResponse
	DecodeTo(&pr, createPriceResp)

	createSessionResp := api.Post("/api/v1/checkout/sessions/", s.CreateStripeCheckoutSessionRequest{
		UserID:     user.ID.String(),
		PriceID:    pr.ID,
		SuccessURL: "https://example.com/success",
		CancelURL:  "https://example.com/cancel",
		Quantity:   1,
	}, auth)
	var createdSession s.StripeCheckoutSessionResponse
	DecodeTo(&createdSession, createSessionResp)

	resp := api.Get("/api/v1/checkout/sessions/", s.GetAllStripeSessionsRequest{Limit: 10}, auth)
	var sessionList []*s.StripeCheckoutSessionResponse
	DecodeTo(&sessionList, resp)

	if len(sessionList) == 0 {
		t.Fatalf("expected to find at least 1 session, but got 0")
	}

	found := false
	for _, sess := range sessionList {
		if sess.ID == createdSession.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected session %s to be in the list, but it was not found", createdSession.ID)
	}
}
