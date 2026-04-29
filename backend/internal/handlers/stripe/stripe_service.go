package stripe

import (
	"context"
	"encoding/json"
	"inside-athletics/internal/models"
	"inside-athletics/internal/utils"
	"os"
	"time"

	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	stripego "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/subscription"
	"github.com/stripe/stripe-go/v82/webhook"
)

type StripeService struct {
	db *gorm.DB
}

func NewStripeService(db *gorm.DB) *StripeService {
	return &StripeService{db: db}
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

	product_params := &stripego.ProductParams{
		Name:        stripego.String(input.Body.Name),
		Description: stripego.String(input.Body.Description),
	}

	stripe_product, err := utils.HandleDBError(product.New(product_params))
	if err != nil {
		return nil, err
	}

	response := &StripeProductResponse{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
		Active:      stripe_product.Active,
		Created:     stripe_product.Created,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) CreateStripePrice(ctx context.Context, input *struct{ Body CreateStripePriceRequest }) (*utils.ResponseBody[StripePriceResponse], error) {
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

	price_params := &stripego.PriceParams{
		Product:    stripego.String(input.Body.Product_ID),
		UnitAmount: stripego.Int64(int64(input.Body.UnitAmount)), // multiply by 100 since stripe does not take floats
		Currency:   stripego.String(string(stripego.CurrencyUSD)),  //hardcoded USD
		Recurring: &stripego.PriceRecurringParams{
			Interval:      stripego.String(string(input.Body.Interval)),
			IntervalCount: stripego.Int64(int64(input.Body.IntervalCount)),
		},
	}

	stripe_price, err := price.New(price_params)
	if err != nil {
		return nil, err
	}

	response := &StripePriceResponse{
		ID:         stripe_price.ID,
		ProductID:  stripe_price.Product.ID,
		UnitAmount: stripe_price.UnitAmount,
		Currency:   string(stripe_price.Currency),
		Active:     stripe_price.Active,
		Created:    stripe_price.Created,
	}

	if stripe_price.Recurring != nil {
		response.Interval = string(stripe_price.Recurring.Interval)
		response.IntervalCount = stripe_price.Recurring.IntervalCount
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) GetStripeProductByID(ctx context.Context, input *GetStripeProductByIDParams) (*utils.ResponseBody[StripeProductResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}

	stripe_product, err := product.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	response := &StripeProductResponse{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
		Active:      stripe_product.Active,
		Created:     stripe_product.Created,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) GetStripePriceByID(ctx context.Context, input *GetStripePriceByIDParams) (*utils.ResponseBody[StripePriceResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("price ID is empty")
	}

	stripePrice, err := price.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	var interval string
	var intervalCount int64

	if stripePrice.Recurring != nil {
		interval = string(stripePrice.Recurring.Interval)
		intervalCount = stripePrice.Recurring.IntervalCount
	}

	response := &StripePriceResponse{
		ID:            stripePrice.ID,
		ProductID:     stripePrice.Product.ID,
		UnitAmount:    stripePrice.UnitAmount,
		Currency:      string(stripePrice.Currency),
		Interval:      interval,
		IntervalCount: intervalCount,
		Active:        stripePrice.Active,
		Created:       stripePrice.Created,
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) UpdateStripeProduct(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body UpdateStripeProductRequest
}) (*utils.ResponseBody[StripeProductResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}

	product_params := &stripego.ProductParams{
		Name:        input.Body.Name,
		Description: input.Body.Description,
	}

	stripe_product, err := product.Update(input.ID, product_params)
	if err != nil {
		return nil, err
	}

	response := &StripeProductResponse{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
		Active:      stripe_product.Active,
		Created:     stripe_product.Created,
	}

	return &utils.ResponseBody[StripeProductResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) UpdateStripePrice(ctx context.Context, input *struct {
	ID   string `path:"id"`
	Body UpdateStripePriceRequest
}) (*utils.ResponseBody[StripePriceResponse], error) {

	if input.ID == "" {
		return nil, fmt.Errorf("price ID is empty")
	}

	oldPrice, err := price.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	newPriceParams := &stripego.PriceParams{
		Product:    stripego.String(oldPrice.Product.ID),
		UnitAmount: stripego.Int64(int64(*input.Body.UnitAmount)),
		Currency:   stripego.String(string(oldPrice.Currency)),
		Recurring: &stripego.PriceRecurringParams{
			Interval:      stripego.String(string(*input.Body.Interval)),
			IntervalCount: stripego.Int64(int64(*input.Body.IntervalCount)),
		},
	}

	newStripePrice, err := price.New(newPriceParams)
	if err != nil {
		return nil, err
	}

	_, _ = price.Update(input.ID, &stripego.PriceParams{
		Active: stripego.Bool(false),
	})

	var interval string
	var intervalCount int64

	if newStripePrice.Recurring != nil {
		interval = string(newStripePrice.Recurring.Interval)
		intervalCount = newStripePrice.Recurring.IntervalCount
	}

	response := &StripePriceResponse{
		ID:            newStripePrice.ID,
		ProductID:     newStripePrice.Product.ID,
		UnitAmount:    newStripePrice.UnitAmount,
		Currency:      string(newStripePrice.Currency),
		Interval:      interval,
		IntervalCount: intervalCount,
		Active:        newStripePrice.Active,
		Created:       newStripePrice.Created,
	}

	return &utils.ResponseBody[StripePriceResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) GetAllStripeProducts(ctx context.Context, input *GetAllStripeProductsRequest) (*utils.ResponseBody[[]*StripeProductResponse], error) {
	params := &stripego.ProductListParams{}
	iter := product.List(params)
	products := make([]*StripeProductResponse, 0)

	for iter.Next() {
		stripe_product := iter.Product()
		response := &StripeProductResponse{
			ID:          stripe_product.ID,
			Name:        stripe_product.Name,
			Description: stripe_product.Description,
			Active:      stripe_product.Active,
			Created:     stripe_product.Created,
		}
		products = append(products, response)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[[]*StripeProductResponse]{
		Body: &products,
	}, nil
}

func (s *StripeService) GetAllStripePrices(ctx context.Context, input *GetAllStripePricesRequest) (*utils.ResponseBody[[]*StripePriceResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}
	params := &stripego.PriceListParams{
		Product: stripego.String(input.ID),
	}
	iter := price.List(params)
	prices := make([]*StripePriceResponse, 0)

	for iter.Next() {
		stripe_price := iter.Price()
		response := &StripePriceResponse{
			ID:         stripe_price.ID,
			ProductID:  stripe_price.Product.ID,
			UnitAmount: stripe_price.UnitAmount,
			Currency:   string(stripe_price.Currency),
			Active:     stripe_price.Active,
			Created:    stripe_price.Created,
		}
		prices = append(prices, response)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[[]*StripePriceResponse]{
		Body: &prices,
	}, nil
}

// apparently you can't delete a product due to historical billing data, but you can archive it
// archiving a product automatically archives all the prices associated with it
func (s *StripeService) ArchiveStripeProduct(ctx context.Context, input *ArchiveStripeProductRequest) (*utils.ResponseBody[*StripeProductResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}
	params := &stripego.ProductParams{
		Active: stripego.Bool(false),
	}
	stripe_product, err := product.Update(input.ID, params)
	if err != nil {
		return nil, err
	}
	priceParams := &stripego.PriceListParams{
		Product: stripego.String(input.ID),
	}
	i := price.List(priceParams)

	for i.Next() {
		p := i.Price()
		_, err = price.Update(p.ID, &stripego.PriceParams{
			Active: stripego.Bool(false),
		})
		if err != nil {
			continue
		}
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	response := &StripeProductResponse{
		ID:          stripe_product.ID,
		Name:        stripe_product.Name,
		Description: stripe_product.Description,
		Active:      stripe_product.Active,
		Created:     stripe_product.Created,
	}
	return &utils.ResponseBody[*StripeProductResponse]{
		Body: &response,
	}, nil
}

func (s *StripeService) ArchiveStripePrice(ctx context.Context, input *ArchiveStripePriceRequest) (*utils.ResponseBody[*StripePriceResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("price ID is empty")
	}
	params := &stripego.PriceParams{
		Active: stripego.Bool(false),
	}

	stripe_price, err := price.Update(input.ID, params)
	if err != nil {
		return nil, err
	}

	response := &StripePriceResponse{
		ID:         stripe_price.ID,
		ProductID:  stripe_price.Product.ID,
		UnitAmount: stripe_price.UnitAmount,
		Currency:   string(stripe_price.Currency),
		Active:     stripe_price.Active,
		Created:    stripe_price.Created,
	}

	return &utils.ResponseBody[*StripePriceResponse]{
		Body: &response,
	}, nil
}

func (s *StripeService) RegisterStripeCustomer(ctx context.Context, input *RegisterStripeCustomerInput) (*utils.ResponseBody[RegisterStripeCustomerResponse], error) {
	params := &stripego.CustomerParams{
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

	params := &stripego.CustomerParams{}
	result, err := customer.Get(id, params)
	if err != nil {
		return nil, err
	}

	return &utils.ResponseBody[GetStripeCustomerResponse]{
		Body: mapStripeCustomerToModel(result),
	}, nil
}

func (s *StripeService) GetStripeCustomerByEmail(ctx context.Context, input *GetStripeCustomerByEmailInput) (*utils.ResponseBody[GetStripeCustomerByEmailResponse], error) {
	email := input.Email

	params := &stripego.CustomerSearchParams{
		SearchParams: stripego.SearchParams{
			Query: "email:'" + email + "'",
		},
	}
	iter := customer.Search(params)
	if err := iter.Err(); err != nil {
		return nil, err
	}

	if !iter.Next() {
		return nil, fmt.Errorf("customer not found with that email")
	}
	cust := iter.Customer()

	return &utils.ResponseBody[GetStripeCustomerByEmailResponse]{
		Body: &GetStripeCustomerByEmailResponse{
			ID:    cust.ID,
			Email: cust.Email,
		},
	}, nil
}

func (s *StripeService) UpdateStripeCustomer(ctx context.Context, input *UpdateStripeCustomerInput) (*utils.ResponseBody[GetStripeCustomerResponse], error) {
	id := input.ID

	params := &stripego.CustomerParams{
		Name:          input.Body.Name,
		Email:         input.Body.Email,
		Phone:         input.Body.Phone,
		Description:   input.Body.Description,
		DefaultSource: input.Body.DefaultSource,
		Balance:       input.Body.Balance,
	}

	if input.Body.TaxExempt != nil {
		params.TaxExempt = stripego.String(*input.Body.TaxExempt)
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

	params := &stripego.CustomerParams{}

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

func mapStripeCustomerToModel(c *stripego.Customer) *GetStripeCustomerResponse {
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

func (s *StripeService) CheckActiveSubscription(ctx context.Context, input *HasActiveSubscriptionInput) (*utils.ResponseBody[HasActiveSubscriptionResponse], error) {
	id := input.CustomerID

	params := &stripego.SubscriptionListParams{
		Customer: stripego.String(id),
		Status:   stripego.String("all"),
	}

	iter := subscription.List(params)

	for iter.Next() {
		sub := iter.Subscription()

		if sub.Status == stripego.SubscriptionStatusActive {
			respBody := &HasActiveSubscriptionResponse{
				HasActiveSubscription: true,
				SubscriptionID:        sub.ID,
				Status:                string(sub.Status),
				CurrentPeriodEnd:      sub.BillingCycleAnchor,
			}
			return &utils.ResponseBody[HasActiveSubscriptionResponse]{
				Body: respBody,
			}, nil
		}
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[HasActiveSubscriptionResponse]{
		Body: &HasActiveSubscriptionResponse{
			HasActiveSubscription: false,
		},
	}, nil
}

func (s *StripeService) CreateStripeCheckoutSession(
	ctx context.Context,
	input *struct {
		Body CreateStripeCheckoutSessionRequest
	},
) (*utils.ResponseBody[StripeCheckoutSessionResponse], error) {

	if input.Body.UserID == "" {
		return nil, huma.Error422UnprocessableEntity("user_id cannot be empty.")
	}
	if input.Body.PriceID == "" {
		return nil, huma.Error422UnprocessableEntity("price_id cannot be empty.")
	}
	if input.Body.SuccessURL == "" {
		return nil, huma.Error422UnprocessableEntity("success_url cannot be empty.")
	}
	if input.Body.CancelURL == "" {
		return nil, huma.Error422UnprocessableEntity("cancel_url cannot be empty.")
	}
	if input.Body.Quantity <= 0 {
		return nil, huma.Error422UnprocessableEntity("quantity must be greater than 0.")
	}

	userID, err := uuid.Parse(input.Body.UserID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("user_id is not a valid UUID.")
	}

	// Look up user and get or create their Stripe customer
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, huma.Error404NotFound("user not found")
	}

	customerID, err := s.getOrCreateStripeCustomer(&user)
	if err != nil {
		return nil, err
	}

	params := &stripego.CheckoutSessionParams{
		Params: stripego.Params{
			Expand: []*string{
				stripego.String("line_items"),
				stripego.String("subscription"),
			},
		},
		Customer:   stripego.String(customerID),
		Mode:       stripego.String(string(stripego.CheckoutSessionModeSubscription)),
		SuccessURL: stripego.String(input.Body.SuccessURL),
		CancelURL:  stripego.String(input.Body.CancelURL),
		LineItems: []*stripego.CheckoutSessionLineItemParams{
			{
				Price:    stripego.String(input.Body.PriceID),
				Quantity: stripego.Int64(input.Body.Quantity),
			},
		},
		SubscriptionData: &stripego.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": userID.String(),
			},
		},
	}

	stripeSession, err := utils.HandleDBError(session.New(params))
	if err != nil {
		return nil, err
	}

	response := &StripeCheckoutSessionResponse{
		ID:      stripeSession.ID,
		URL:     stripeSession.URL,
		Mode:    string(stripeSession.Mode),
		Status:  string(stripeSession.Status),
		Created: stripeSession.Created,
	}

	return &utils.ResponseBody[StripeCheckoutSessionResponse]{
		Body: response,
	}, nil
}

// getOrCreateStripeCustomer returns the user's Stripe customer ID, creating one if needed.
func (s *StripeService) getOrCreateStripeCustomer(user *models.User) (string, error) {
	if user.StripeCustomerID != nil && *user.StripeCustomerID != "" {
		return *user.StripeCustomerID, nil
	}

	params := &stripego.CustomerParams{
		Email: stripego.String(user.Email),
		Name:  stripego.String(user.FirstName + " " + user.LastName),
		Metadata: map[string]string{
			"user_id": user.ID.String(),
		},
	}
	cust, err := customer.New(params)
	if err != nil {
		return "", err
	}

	if err := s.db.Model(user).Update("stripe_customer_id", cust.ID).Error; err != nil {
		return "", err
	}
	return cust.ID, nil
}

// HandleWebhook is a raw Fiber handler (not Huma) that processes Stripe webhook events.
func (s *StripeService) HandleWebhook(c *fiber.Ctx) error {
	payload := c.Body()
	sigHeader := c.Get("Stripe-Signature")
	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	event, err := webhook.ConstructEventWithOptions(payload, sigHeader, secret, webhook.ConstructEventOptions{
		IgnoreAPIVersionMismatch: true,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid webhook signature: " + err.Error())
	}

	switch event.Type {
	case "checkout.session.completed":
		var sess stripego.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("failed to parse session")
		}
		if err := s.handleCheckoutCompleted(&sess); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

	case "customer.subscription.updated":
		var sub stripego.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("failed to parse subscription")
		}
		if err := s.handleSubscriptionUpdated(&sub); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

	case "customer.subscription.deleted":
		var sub stripego.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("failed to parse subscription")
		}
		if err := s.handleSubscriptionDeleted(&sub); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

	case "invoice.payment_failed":
		var inv stripego.Invoice
		if err := json.Unmarshal(event.Data.Raw, &inv); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("failed to parse invoice")
		}
		if err := s.handlePaymentFailed(&inv); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *StripeService) handleCheckoutCompleted(sess *stripego.CheckoutSession) error {
	fmt.Printf("[checkout] customer=%v subscription=%v\n", sess.Customer, sess.Subscription)
	if sess.Customer == nil || sess.Subscription == nil {
		fmt.Println("[checkout] early return: customer or subscription is nil")
		return nil
	}

	var user models.User
	if err := s.db.First(&user, "stripe_customer_id = ?", sess.Customer.ID).Error; err != nil {
		fmt.Printf("[checkout] user not found for customer %s: %v\n", sess.Customer.ID, err)
		return fmt.Errorf("user not found for customer %s", sess.Customer.ID)
	}
	fmt.Printf("[checkout] found user %s, granting premium role\n", user.ID)

	// Fetch full subscription from Stripe since webhook payload doesn't expand it
	fullSub, err := subscription.Get(sess.Subscription.ID, nil)
	if err != nil {
		return fmt.Errorf("failed to fetch subscription %s: %w", sess.Subscription.ID, err)
	}

	priceID := ""
	if fullSub.Items != nil && len(fullSub.Items.Data) > 0 {
		priceID = fullSub.Items.Data[0].Price.ID
	}

	now := time.Now()
	record := models.UserSubscription{
		UserID:               user.ID,
		StripeSubscriptionID: fullSub.ID,
		StripePriceID:        priceID,
		Status:               models.SubscriptionStatusActive,
		CurrentPeriodStart:   now,
		CurrentPeriodEnd:     now.AddDate(0, 1, 0),
	}

	if err := s.db.Where(models.UserSubscription{UserID: user.ID}).
		Assign(record).
		FirstOrCreate(&record).Error; err != nil {
		return err
	}

	fmt.Printf("[checkout] calling grantPremiumRole for user %s\n", user.ID)
	err = s.grantPremiumRole(user.ID)
	fmt.Printf("[checkout] grantPremiumRole result: %v\n", err)
	return err
}

func (s *StripeService) handleSubscriptionUpdated(sub *stripego.Subscription) error {
	if sub.Customer == nil {
		return nil
	}

	var user models.User
	if err := s.db.First(&user, "stripe_customer_id = ?", sub.Customer.ID).Error; err != nil {
		return nil // user not found — ignore
	}

	status := models.SubscriptionStatus(sub.Status)

	updates := map[string]interface{}{
		"stripe_subscription_id": sub.ID,
		"status":                 status,
	}

	return s.db.Model(&models.UserSubscription{}).
		Where("user_id = ?", user.ID).
		Updates(updates).Error
}

func (s *StripeService) handleSubscriptionDeleted(sub *stripego.Subscription) error {
	if sub.Customer == nil {
		return nil
	}

	var user models.User
	if err := s.db.First(&user, "stripe_customer_id = ?", sub.Customer.ID).Error; err != nil {
		return nil
	}

	now := time.Now()
	if err := s.db.Model(&models.UserSubscription{}).
		Where("user_id = ?", user.ID).
		Updates(map[string]interface{}{
			"status":      models.SubscriptionStatusCanceled,
			"canceled_at": now,
		}).Error; err != nil {
		return err
	}

	return s.revokePremiumRole(user.ID)
}

func (s *StripeService) handlePaymentFailed(inv *stripego.Invoice) error {
	if inv.Customer == nil {
		return nil
	}

	var user models.User
	if err := s.db.First(&user, "stripe_customer_id = ?", inv.Customer.ID).Error; err != nil {
		return nil
	}

	return s.db.Model(&models.UserSubscription{}).
		Where("user_id = ?", user.ID).
		Update("status", models.SubscriptionStatusPastDue).Error
}

func (s *StripeService) grantPremiumRole(userID uuid.UUID) error {
	var premiumRole models.Role
	if err := s.db.Where("name = ?", models.RolePremiumUser).First(&premiumRole).Error; err != nil {
		return fmt.Errorf("premium_user role not found: %w", err)
	}
	// Upsert: replaces any existing role (e.g. "user") with premium_user.
	// Admins/moderators already satisfy hasPremium so we skip them.
	return s.db.Exec(`
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?)
		ON CONFLICT (user_id) DO UPDATE SET role_id = EXCLUDED.role_id
		WHERE user_roles.role_id NOT IN (
			SELECT id FROM roles WHERE name IN ('admin', 'moderator')
		)`,
		userID, premiumRole.ID,
	).Error
}

func (s *StripeService) revokePremiumRole(userID uuid.UUID) error {
	var premiumRole models.Role
	if err := s.db.Where("name = ?", models.RolePremiumUser).First(&premiumRole).Error; err != nil {
		return nil
	}
	var userRole models.Role
	if err := s.db.Where("name = ?", models.RoleUser).First(&userRole).Error; err != nil {
		return fmt.Errorf("user role not found: %w", err)
	}
	// Only downgrade if they specifically have the premium_user role.
	return s.db.Exec(
		"UPDATE user_roles SET role_id = ? WHERE user_id = ? AND role_id = ?",
		userRole.ID, userID, premiumRole.ID,
	).Error
}

func (s *StripeService) GetStripeCheckoutSessionByID(ctx context.Context, input *GetStripeCheckoutSessionParams) (*utils.ResponseBody[StripeCheckoutSessionResponse], error) {
	if input.ID == "" {
		return nil, fmt.Errorf("product ID is empty")
	}

	stripeSession, err := session.Get(input.ID, nil)
	if err != nil {
		return nil, err
	}

	response := &StripeCheckoutSessionResponse{
		ID:      stripeSession.ID,
		URL:     stripeSession.URL,
		Mode:    string(stripeSession.Mode),
		Status:  string(stripeSession.Status),
		Created: stripeSession.Created,
	}

	return &utils.ResponseBody[StripeCheckoutSessionResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) DeleteStripeCheckoutSession(
	ctx context.Context, input *DeleteCheckoutSessionRequest,
) (*utils.ResponseBody[StripeCheckoutSessionResponse], error) {

	if input.ID == "" {
		return nil, huma.Error422UnprocessableEntity("session id cannot be empty")
	}

	stripeSession, err := session.Expire(input.ID, nil)
	if err != nil {
		return nil, err
	}

	response := &StripeCheckoutSessionResponse{
		ID:      stripeSession.ID,
		URL:     stripeSession.URL,
		Mode:    string(stripeSession.Mode),
		Status:  string(stripeSession.Status),
		Created: stripeSession.Created,
	}

	return &utils.ResponseBody[StripeCheckoutSessionResponse]{
		Body: response,
	}, nil
}

func (s *StripeService) GetAllStripeSessions(
	ctx context.Context, input *GetAllStripeSessionsRequest,
) (*utils.ResponseBody[[]*StripeCheckoutSessionResponse], error) {

	limit := int64(50)
	if input.Limit > 0 {
		limit = input.Limit
	}

	params := &stripego.CheckoutSessionListParams{
		ListParams: stripego.ListParams{
			Limit: &limit,
		},
	}

	if input.CustomerID != "" {
		params.Customer = stripego.String(input.CustomerID)
	}

	i := session.List(params)
	var sessions []*StripeCheckoutSessionResponse

	for i.Next() {
		stripeSession := i.CheckoutSession()
		response := &StripeCheckoutSessionResponse{
			ID:      stripeSession.ID,
			URL:     stripeSession.URL,
			Mode:    string(stripeSession.Mode),
			Status:  string(stripeSession.Status),
			Created: stripeSession.Created,
		}
		sessions = append(sessions, response)
	}

	if err := i.Err(); err != nil {
		return nil, err
	}

	return &utils.ResponseBody[[]*StripeCheckoutSessionResponse]{
		Body: &sessions,
	}, nil
}
