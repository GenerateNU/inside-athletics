package routeTests

import (
	s "inside-athletics/internal/handlers/stripe"
	"inside-athletics/internal/utils"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/product"
	// "github.com/stripe/stripe-go/v72/price"
)

func TestCreateProduct(t *testing.T) {
	testDB := SetupTestDB(t)
	defer testDB.Teardown(t)
	api := testDB.API
	t.Setenv("STRIPE_TEST_KEY", "sk_test_51SyjYFLGVwetm7oJsQ1yKE7vYFJoQxAXNoGqhgrIRcCpjYuMZbVwPkXsuZnfMmNgyDRaE32bAMFDYhXiHRuunYBd00LFwWupaT")
	stripe.Key = os.Getenv("STRIPE_TEST_KEY")

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	body := s.CreateStripeProductRequest{
		Name:        name,
		Description: description,
	}

	resp := api.Post("/api/v1/stripe_product/", body, "Authorization: Bearer "+uuid.NewString())
	var product stripe.Product
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
	t.Setenv("STRIPE_TEST_KEY", "sk_test_51SyjYFLGVwetm7oJsQ1yKE7vYFJoQxAXNoGqhgrIRcCpjYuMZbVwPkXsuZnfMmNgyDRaE32bAMFDYhXiHRuunYBd00LFwWupaT")
	stripe.Key = os.Getenv("STRIPE_TEST_KEY")

	name := "Premium Plan"
	description := "Get premium content with this subscription"

	product_params := &stripe.ProductParams{
		Name:        stripe.String(name),
		Description: stripe.String(description),
	}

	resp, _ := utils.HandleDBError(product.New(product_params))

	body := s.GetStripeProductByIDParams{
		ID: resp.ID,
	}

	resp_recieved := api.Get("/api/v1/stripe_product/"+resp.ID, body, "Authorization: Bearer "+uuid.NewString())
	var product stripe.Product
	DecodeTo(&product, resp_recieved)

	if product.Name != "Premium Plan" {
		t.Errorf("expected Name to be 'Premium Plan', got %s", product.Name)
	}
	if product.Description != "Get premium content with this subscription" {
		t.Errorf("expected Description to be 'Get premium content...', got %s", product.Description)
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

	t.Logf(("running endpoint"))
	resp := api.Post("/api/v1/stripe_customers", reqBody, "Authorization: Bearer "+uuid.NewString())

	t.Logf("Response code: %d", resp.Code)
	bodyStr := resp.Body.String()
	t.Logf("Response body: %s", bodyStr)

	// t.Logf("Response status: %d", resp.StatusCode)
	// t.Logf("Response headers: %+v", resp.Header)

	var c s.RegisterStripeCustomerResponse
	DecodeTo(&c, resp)

	t.Logf("Decoded response: %+v", c)

	if c.ID == "" {
		t.Fatalf("Expected ID to be set, got nil")
	}

	// params := &stripe.CustomerParams{}
	// customer, err := customer.Get(c.ID, params)

	// if err != nil {
	// 	t.Fatalf("Unexpected response: %+v", err)
	// }
	// if customer != nil {
	// 	t.Fatalf("Unexpected")
	// }

	// if customer.Name != name ||
	// 	customer.Email != email ||
	// 	customer.Phone != phone ||
	// 	customer.Description != description {
	// 	t.Fatalf("Unexpected response: %+v", c)
	// }

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
