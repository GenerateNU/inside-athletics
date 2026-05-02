package stripe

import (
	stripego "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/subscription"
)

// StripeClient abstracts all Stripe API calls so the service can be tested with a mock.
type StripeClient interface {
	// Products
	CreateProduct(params *stripego.ProductParams) (*stripego.Product, error)
	GetProduct(id string) (*stripego.Product, error)
	UpdateProduct(id string, params *stripego.ProductParams) (*stripego.Product, error)
	ListProducts(params *stripego.ProductListParams) ([]*stripego.Product, error)

	// Prices
	CreatePrice(params *stripego.PriceParams) (*stripego.Price, error)
	GetPrice(id string) (*stripego.Price, error)
	UpdatePrice(id string, params *stripego.PriceParams) (*stripego.Price, error)
	ListPrices(params *stripego.PriceListParams) ([]*stripego.Price, error)

	// Customers
	CreateCustomer(params *stripego.CustomerParams) (*stripego.Customer, error)
	GetCustomer(id string) (*stripego.Customer, error)
	SearchCustomers(params *stripego.CustomerSearchParams) ([]*stripego.Customer, error)
	UpdateCustomer(id string, params *stripego.CustomerParams) (*stripego.Customer, error)
	DeleteCustomer(id string) (*stripego.Customer, error)

	// Checkout sessions
	CreateSession(params *stripego.CheckoutSessionParams) (*stripego.CheckoutSession, error)
	GetSession(id string) (*stripego.CheckoutSession, error)
	ExpireSession(id string) (*stripego.CheckoutSession, error)
	ListSessions(params *stripego.CheckoutSessionListParams) ([]*stripego.CheckoutSession, error)

	// Subscriptions
	GetSubscription(id string) (*stripego.Subscription, error)
	ListSubscriptions(params *stripego.SubscriptionListParams) ([]*stripego.Subscription, error)
}

// realStripeClient calls the live Stripe API.
type realStripeClient struct{}

func (c *realStripeClient) CreateProduct(p *stripego.ProductParams) (*stripego.Product, error) {
	return product.New(p)
}
func (c *realStripeClient) GetProduct(id string) (*stripego.Product, error) {
	return product.Get(id, nil)
}
func (c *realStripeClient) UpdateProduct(id string, p *stripego.ProductParams) (*stripego.Product, error) {
	return product.Update(id, p)
}
func (c *realStripeClient) ListProducts(p *stripego.ProductListParams) ([]*stripego.Product, error) {
	var out []*stripego.Product
	iter := product.List(p)
	for iter.Next() {
		out = append(out, iter.Product())
	}
	return out, iter.Err()
}

func (c *realStripeClient) CreatePrice(p *stripego.PriceParams) (*stripego.Price, error) {
	return price.New(p)
}
func (c *realStripeClient) GetPrice(id string) (*stripego.Price, error) {
	return price.Get(id, nil)
}
func (c *realStripeClient) UpdatePrice(id string, p *stripego.PriceParams) (*stripego.Price, error) {
	return price.Update(id, p)
}
func (c *realStripeClient) ListPrices(p *stripego.PriceListParams) ([]*stripego.Price, error) {
	var out []*stripego.Price
	iter := price.List(p)
	for iter.Next() {
		out = append(out, iter.Price())
	}
	return out, iter.Err()
}

func (c *realStripeClient) CreateCustomer(p *stripego.CustomerParams) (*stripego.Customer, error) {
	return customer.New(p)
}
func (c *realStripeClient) GetCustomer(id string) (*stripego.Customer, error) {
	return customer.Get(id, nil)
}
func (c *realStripeClient) SearchCustomers(p *stripego.CustomerSearchParams) ([]*stripego.Customer, error) {
	var out []*stripego.Customer
	iter := customer.Search(p)
	for iter.Next() {
		out = append(out, iter.Customer())
	}
	return out, iter.Err()
}
func (c *realStripeClient) UpdateCustomer(id string, p *stripego.CustomerParams) (*stripego.Customer, error) {
	return customer.Update(id, p)
}
func (c *realStripeClient) DeleteCustomer(id string) (*stripego.Customer, error) {
	result, err := customer.Del(id, nil)
	if err != nil {
		return nil, err
	}
	return &stripego.Customer{ID: result.ID}, nil
}

func (c *realStripeClient) CreateSession(p *stripego.CheckoutSessionParams) (*stripego.CheckoutSession, error) {
	return session.New(p)
}
func (c *realStripeClient) GetSession(id string) (*stripego.CheckoutSession, error) {
	return session.Get(id, nil)
}
func (c *realStripeClient) ExpireSession(id string) (*stripego.CheckoutSession, error) {
	return session.Expire(id, nil)
}
func (c *realStripeClient) ListSessions(p *stripego.CheckoutSessionListParams) ([]*stripego.CheckoutSession, error) {
	var out []*stripego.CheckoutSession
	iter := session.List(p)
	for iter.Next() {
		out = append(out, iter.CheckoutSession())
	}
	return out, iter.Err()
}

func (c *realStripeClient) GetSubscription(id string) (*stripego.Subscription, error) {
	return subscription.Get(id, nil)
}
func (c *realStripeClient) ListSubscriptions(p *stripego.SubscriptionListParams) ([]*stripego.Subscription, error) {
	var out []*stripego.Subscription
	iter := subscription.List(p)
	for iter.Next() {
		out = append(out, iter.Subscription())
	}
	return out, iter.Err()
}
