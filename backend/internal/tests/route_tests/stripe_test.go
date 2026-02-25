package routeTests

import (
	"fmt"
	s "inside-athletics/internal/handlers/stripe"
	"testing"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/customer"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
)

func TestCreateProduct(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	body := s.CreateStripeProductRequest{
		Name:        name,
		Description: description,
	}

	resp := api.Post("/api/v1/stripe_product/", body, "Authorization: Bearer "+uuid.NewString())
	var product s.StripeProductResponse
	DecodeTo(&product, resp)

	if product.Name != "Premium Plan" {
		t.Errorf("expected Name to be 'Premium Plan', got %s", product.Name)
	}
	if product.Description != "Get premium content with this subscription" {
		t.Errorf("expected Description to be 'Get premium content...', got %s", product.Description)
	}
}

func TestGetProductByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	result, err := product.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID

	body := s.GetStripeProductByIDParams{
		ID: id,
	}

	resp := api.Get(
		"/api/v1/stripe_product/"+id, body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var stripeProduct s.StripeProductResponse
	DecodeTo(&stripeProduct, resp)

	if stripeProduct.Name != name {
		t.Errorf("expected Name to be '%s', got %s", name, stripeProduct.Name)
	}

	if stripeProduct.Description != description {
		t.Errorf("expected Description to be '%s', got %s", description, stripeProduct.Description)
	}
}

func TestUpdateStripeProduct(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	result, err := product.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID

	newName := "Basic Plan"
	newDescription := "The plan you get when you wanna miss out on premium content..."

	body := s.UpdateStripeProductRequest{
		Name:        stripe.String(newName),
		Description: stripe.String(newDescription),
	}

	resp := api.Patch(
		"/api/v1/stripe_product/"+id, body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var stripeProduct s.StripeProductResponse
	DecodeTo(&stripeProduct, resp)

	if stripeProduct.Name != newName {
		t.Errorf("expected Name to be '%s', got %s", name, stripeProduct.Name)
	}

	if stripeProduct.Description != newDescription {
		t.Errorf("expected Description to be '%s', got %s", description, stripeProduct.Description)
	}
}

func TestArchiveStripeProduct(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	result, err := product.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID

	body := s.ArchiveStripeProductRequest{
		ID: id,
	}

	resp := api.Delete(
		"/api/v1/stripe_product/"+id, body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var stripeProduct s.StripeProductResponse
	DecodeTo(&stripeProduct, resp)

	if stripeProduct.Active {
		t.Errorf("expected product to be inactive after delete, but Active = true")
	}
}

func TestGetAllProducts(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	_, err := product.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	name = "Basic Plan"
	description = "Get basic content with this subscription"

	params = &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	_, err = product.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	body := s.GetAllStripeProductsRequest{}

	resp := api.Get(
		"/api/v1/stripe_products/", body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var stripeProducts []s.StripeProductResponse
	DecodeTo(&stripeProducts, resp)

	if len(stripeProducts) < 2 {
		t.Errorf("expected at least 2 products, got %d", len(stripeProducts))
	}
}

func TestCreatePrice(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	product_params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	result, err := product.New(product_params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID
	unitAmount := 2500
	interval := s.Day
	intervalCount := 3

	body := s.CreateStripePriceRequest{
		Product_ID:    id,
		UnitAmount:    unitAmount,
		Interval:      interval,
		IntervalCount: intervalCount,
	}

	resp := api.Post("/api/v1/stripe_price/", body, "Authorization: Bearer "+uuid.NewString())
	var price s.StripePriceResponse
	DecodeTo(&price, resp)

	if price.ID == "" {
		t.Errorf("expected price ID to be set, got empty string")
	}

	if price.ProductID != id {
		t.Errorf("expected product id to be %s, got %s", id, price.ProductID)
	}

	if price.UnitAmount != int64(unitAmount) {
		t.Errorf("expected unit amount to be %d, got %d", unitAmount, price.UnitAmount)
	}

	if price.Interval != string(interval) {
		t.Errorf("expected interval to be %s, got %s", interval, price.Interval)
	}

	if price.IntervalCount != int64(intervalCount) {
		t.Errorf("expected interval count to be %d, got %d", intervalCount, price.IntervalCount)
	}

}

func TestGetStripePriceByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	product_params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	result, err := product.New(product_params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	id := result.ID
	unitAmount := 2550
	interval := s.Day
	intervalCount := 3

	price_params := &stripe.PriceParams{
		Product:    stripe.String(id),
		UnitAmount: stripe.Int64(int64(unitAmount)),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),

		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(interval)),
			IntervalCount: stripe.Int64(int64(intervalCount)),
		},
	}

	price_result, err := price.New(price_params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	price_id := price_result.ID

	body := s.GetStripePriceByIDParams{
		ID: price_id,
	}

	resp := api.Get("/api/v1/stripe_price/"+price_id, body, "Authorization: Bearer "+uuid.NewString())
	var price s.StripePriceResponse
	DecodeTo(&price, resp)

	if price.ID == "" {
		t.Errorf("expected price ID to be set, got empty string")
	}

	if price.ProductID != id {
		t.Errorf("expected product id to be %s, got %s", id, price.ProductID)
	}

	if price.UnitAmount != int64(unitAmount) {
		t.Errorf("expected unit amount to be %d, got %d", unitAmount, price.UnitAmount)
	}

	if price.Interval != string(interval) {
		t.Errorf("expected interval to be %s, got %s", interval, price.Interval)
	}

	if price.IntervalCount != int64(intervalCount) {
		t.Errorf("expected interval count to be %d, got %d", intervalCount, price.IntervalCount)
	}
}

func TestUpdateStripePrice(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	productParams := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	productResult, err := product.New(productParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	productID := productResult.ID

	originalUnitAmount := 2550
	originalInterval := s.Day
	originalIntervalCount := 3

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(int64(originalUnitAmount)),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(originalInterval)),
			IntervalCount: stripe.Int64(int64(originalIntervalCount)),
		},
	}

	originalPrice, err := price.New(priceParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	newUnitAmount := 3000
	newInterval := s.Week
	newIntervalCount := 2

	body := s.UpdateStripePriceRequest{
		UnitAmount:    &newUnitAmount,
		Interval:      &newInterval,
		IntervalCount: &newIntervalCount,
	}

	resp := api.Patch("/api/v1/stripe_price/"+originalPrice.ID, body, "Authorization: Bearer "+uuid.NewString())

	var updatedPrice s.StripePriceResponse
	DecodeTo(&updatedPrice, resp)

	if updatedPrice.ID == originalPrice.ID {
		t.Errorf("expected new price ID, but got same ID %s", updatedPrice.ID)
	}

	if updatedPrice.ProductID != productID {
		t.Errorf("expected product id to be %s, got %s",
			productID, updatedPrice.ProductID)
	}

	if updatedPrice.UnitAmount != int64(newUnitAmount) {
		t.Errorf("expected unit amount to be %d, got %d",
			newUnitAmount, updatedPrice.UnitAmount)
	}

	if updatedPrice.Interval != string(newInterval) {
		t.Errorf("expected interval to be %s, got %s",
			newInterval, updatedPrice.Interval)
	}

	if updatedPrice.IntervalCount != int64(newIntervalCount) {
		t.Errorf("expected interval count to be %d, got %d",
			newIntervalCount, updatedPrice.IntervalCount)
	}
}

func TestArchiveStripePrice(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	productParams := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	productResult, err := product.New(productParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	productID := productResult.ID

	originalUnitAmount := 2550
	originalInterval := s.Day
	originalIntervalCount := 3

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(int64(originalUnitAmount)),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(originalInterval)),
			IntervalCount: stripe.Int64(int64(originalIntervalCount)),
		},
	}

	originalPrice, err := price.New(priceParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	body := s.ArchiveStripePriceRequest{
		ID: originalPrice.ID,
	}

	resp := api.Delete(
		"/api/v1/stripe_price/"+originalPrice.ID, body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var archivedPrice s.StripePriceResponse
	DecodeTo(&archivedPrice, resp)

	if archivedPrice.Active {
		t.Errorf("expected price to be inactive after delete, but Active = true")
	}
}

func TestGetAllStripePrices(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	productParams := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	productResult, err := product.New(productParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	productID := productResult.ID

	originalUnitAmount := 2550
	originalInterval := s.Day
	originalIntervalCount := 3

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(int64(originalUnitAmount)),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(originalInterval)),
			IntervalCount: stripe.Int64(int64(originalIntervalCount)),
		},
	}

	_, err = price.New(priceParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	priceParams = &stripe.PriceParams{
		Product:    stripe.String(productID),
		UnitAmount: stripe.Int64(int64(originalUnitAmount)),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(originalInterval)),
			IntervalCount: stripe.Int64(int64(originalIntervalCount)),
		},
	}

	_, err = price.New(priceParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	body := s.GetAllStripePricesRequest{
		ID: productID,
	}

	resp := api.Get("/api/v1/stripe_prices/"+productID, body, "Authorization: Bearer "+uuid.NewString())

	fmt.Println(resp.Body.String())

	var stripePrices []s.StripePriceResponse
	DecodeTo(&stripePrices, resp)

	if len(stripePrices) < 2 {
		t.Errorf("expected at least 2 products, got %d", len(stripePrices))
	}
}

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
func TestGetCustomerByEmail(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	name := "Suli"
	email := "suli_newemail@gmail.com"
	phone := "888 420 6769"
	description := "premium content user"
	params := &stripe.CustomerParams{
		Name:        &name,
		Email:       &email,
		Phone:       &phone,
		Description: &description,
	}
	_, err := customer.New(params)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	// result_id := result.ID

	resp := api.Get("/api/v1/stripe_customers/email/"+email, "Authorization: Bearer "+uuid.NewString())

	var c s.GetStripeCustomerByEmailResponse

	DecodeTo(&c, resp)

	if c.Email != email {
		t.Fatalf("Unexpected response, email does not match: %+v", c)
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

func TestCreateCheckoutSession(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	productParams := &stripe.ProductParams{
		Name:        stripe.String("Premium Plan"),
		Description: stripe.String("Get premium content with this subscription"),
	}

	productResult, err := product.New(productParams)
	if err != nil {
		t.Fatalf("Failed to create product: %+v", err)
	}

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(productResult.ID),
		UnitAmount: stripe.Int64(2550),
		Currency:   stripe.String("usd"),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(stripe.PriceRecurringIntervalDay)),
			IntervalCount: stripe.Int64(3),
		},
	}

	priceResult, err := price.New(priceParams)
	if err != nil {
		t.Fatalf("Failed to create price: %+v", err)
	}

	reqBody := s.CreateStripeCheckoutSessionRequest{
		PriceID:    priceResult.ID,
		SuccessURL: "https://example.com/success",
		CancelURL:  "https://example.com/cancel",
		Quantity:   2,
	}

	resp := api.Post("/api/v1/checkout/sessions/", reqBody, "Authorization: Bearer "+uuid.NewString())

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

	if string(session.Mode) != string(stripe.CheckoutSessionModeSubscription) {
		t.Errorf("expected mode to be subscription, got %s", session.Mode)
	}
}

func TestGetStripeCheckoutSessionByID(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	productParams := &stripe.ProductParams{
		Name:        stripe.String("Premium Plan"),
		Description: stripe.String("Get premium content with this subscription"),
	}
	productResult, err := product.New(productParams)
	if err != nil {
		t.Fatalf("Failed to create product: %+v", err)
	}

	priceParams := &stripe.PriceParams{
		Product:    stripe.String(productResult.ID),
		UnitAmount: stripe.Int64(2550),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
		Recurring: &stripe.PriceRecurringParams{
			Interval:      stripe.String(string(stripe.PriceRecurringIntervalDay)),
			IntervalCount: stripe.Int64(3),
		},
	}
	priceResult, err := price.New(priceParams)
	if err != nil {
		t.Fatalf("Failed to create price: %+v", err)
	}

	sessionParams := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String("https://example.com/success"),
		CancelURL:  stripe.String("https://example.com/cancel"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceResult.ID),
				Quantity: stripe.Int64(2),
			},
		},
	}
	createdSession, err := session.New(sessionParams)
	if err != nil {
		t.Fatalf("Failed to create checkout session: %+v", err)
	}

	body := s.GetStripeCheckoutSessionParams{
		ID: createdSession.ID,
	}

	getResp := api.Get(
		"/api/v1/checkout/sessions/"+createdSession.ID,
		body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var session s.StripeCheckoutSessionResponse
	DecodeTo(&session, getResp)

	if session.ID == "" {
		t.Errorf("expected checkout session ID to be set, got empty string")
	}

	if session.URL == "" {
		t.Errorf("expected checkout session URL to be set, got empty string")
	}

	if session.Mode != "subscription" {
		t.Errorf("expected mode to be subscription, got %s", session.Mode)
	}

	if string(session.Mode) != string(stripe.CheckoutSessionModeSubscription) {
		t.Errorf("expected mode to be subscription, got %s", session.Mode)
	}
}

func TestDeleteStripeCheckoutSession(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	sessionParams := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Test Product"),
					},
					UnitAmount: stripe.Int64(1000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://example.com/success"),
		CancelURL:  stripe.String("https://example.com/cancel"),
	}

	createdSession, err := session.New(sessionParams)
	if err != nil {
		t.Fatalf("Unexpected response: %+v", err)
	}

	body := s.DeleteCheckoutSessionRequest{
		ID: createdSession.ID,
	}

	resp := api.Delete(
		"/api/v1/checkout/sessions/"+createdSession.ID,
		body,
		"Authorization: Bearer "+uuid.NewString(),
	)

	var session s.StripeCheckoutSessionResponse
	DecodeTo(&session, resp)

	if session.ID == "" {
		t.Errorf("expected checkout session ID to be set, got empty string")
	}
}

func TestGetAllStripeSessions(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API

	p, err := product.New(&stripe.ProductParams{
		Name: stripe.String("Standard Plan"),
	})
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	pr, err := price.New(&stripe.PriceParams{
		Product:    stripe.String(p.ID),
		UnitAmount: stripe.Int64(1000),
		Currency:   stripe.String(string(stripe.CurrencyUSD)),
	})
	if err != nil {
		t.Fatalf("Failed to create price: %v", err)
	}

	sessionParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://example.com/success"),
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(pr.ID),
				Quantity: stripe.Int64(1),
			},
		},
	}
	createdSession, err := session.New(sessionParams)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	requestBody := s.GetAllStripeSessionsRequest{
		Limit: 10,
	}

	resp := api.Get("/api/v1/checkout/sessions/", requestBody, "Authorization: Bearer "+uuid.NewString())
	var sessionList []*s.StripeCheckoutSessionResponse
	DecodeTo(&sessionList, resp)

	if len(sessionList) == 0 {
		t.Fatalf("expected to find at least 1 session for price %s, but got 0", pr.ID)
	}

	found := false
	for _, sess := range sessionList {
		if sess.ID == createdSession.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected session %s to be in the filtered list, but it was not found", createdSession.ID)
	}
}
