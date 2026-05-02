package stripe

import (
	"fmt"
	"time"

	stripego "github.com/stripe/stripe-go/v82"
)

// MockStripeClient is an in-memory Stripe client for use in tests.
type MockStripeClient struct {
	products      map[string]*stripego.Product
	prices        map[string]*stripego.Price
	customers     map[string]*stripego.Customer
	sessions      map[string]*stripego.CheckoutSession
	subscriptions map[string]*stripego.Subscription
	counter       int
}

func NewMockStripeClient() *MockStripeClient {
	return &MockStripeClient{
		products:      make(map[string]*stripego.Product),
		prices:        make(map[string]*stripego.Price),
		customers:     make(map[string]*stripego.Customer),
		sessions:      make(map[string]*stripego.CheckoutSession),
		subscriptions: make(map[string]*stripego.Subscription),
	}
}

func (m *MockStripeClient) nextID(prefix string) string {
	m.counter++
	return fmt.Sprintf("%s_mock_%d", prefix, m.counter)
}

// Products

func (m *MockStripeClient) CreateProduct(p *stripego.ProductParams) (*stripego.Product, error) {
	prod := &stripego.Product{
		ID:      m.nextID("prod"),
		Active:  true,
		Created: time.Now().Unix(),
	}
	if p.Name != nil {
		prod.Name = *p.Name
	}
	if p.Description != nil {
		prod.Description = *p.Description
	}
	m.products[prod.ID] = prod
	return prod, nil
}

func (m *MockStripeClient) GetProduct(id string) (*stripego.Product, error) {
	if p, ok := m.products[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("no such product: %s", id)
}

func (m *MockStripeClient) UpdateProduct(id string, p *stripego.ProductParams) (*stripego.Product, error) {
	prod, err := m.GetProduct(id)
	if err != nil {
		return nil, err
	}
	if p.Name != nil {
		prod.Name = *p.Name
	}
	if p.Description != nil {
		prod.Description = *p.Description
	}
	if p.Active != nil {
		prod.Active = *p.Active
	}
	return prod, nil
}

func (m *MockStripeClient) ListProducts(_ *stripego.ProductListParams) ([]*stripego.Product, error) {
	out := make([]*stripego.Product, 0, len(m.products))
	for _, p := range m.products {
		out = append(out, p)
	}
	return out, nil
}

// Prices

func (m *MockStripeClient) CreatePrice(p *stripego.PriceParams) (*stripego.Price, error) {
	pr := &stripego.Price{
		ID:      m.nextID("price"),
		Active:  true,
		Created: time.Now().Unix(),
	}
	if p.UnitAmount != nil {
		pr.UnitAmount = *p.UnitAmount
	}
	if p.Currency != nil {
		pr.Currency = stripego.Currency(*p.Currency)
	}
	if p.Product != nil {
		pr.Product = &stripego.Product{ID: *p.Product}
	}
	if p.Recurring != nil {
		pr.Recurring = &stripego.PriceRecurring{}
		if p.Recurring.Interval != nil {
			pr.Recurring.Interval = stripego.PriceRecurringInterval(*p.Recurring.Interval)
		}
		if p.Recurring.IntervalCount != nil {
			pr.Recurring.IntervalCount = *p.Recurring.IntervalCount
		}
	}
	m.prices[pr.ID] = pr
	return pr, nil
}

func (m *MockStripeClient) GetPrice(id string) (*stripego.Price, error) {
	if p, ok := m.prices[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("no such price: %s", id)
}

func (m *MockStripeClient) UpdatePrice(id string, p *stripego.PriceParams) (*stripego.Price, error) {
	pr, err := m.GetPrice(id)
	if err != nil {
		return nil, err
	}
	if p.Active != nil {
		pr.Active = *p.Active
	}
	return pr, nil
}

func (m *MockStripeClient) ListPrices(p *stripego.PriceListParams) ([]*stripego.Price, error) {
	out := make([]*stripego.Price, 0, len(m.prices))
	for _, pr := range m.prices {
		if p.Product != nil && pr.Product != nil && pr.Product.ID != *p.Product {
			continue
		}
		out = append(out, pr)
	}
	return out, nil
}

// Customers

func (m *MockStripeClient) CreateCustomer(p *stripego.CustomerParams) (*stripego.Customer, error) {
	c := &stripego.Customer{
		ID:      m.nextID("cus"),
		Created: time.Now().Unix(),
	}
	if p.Name != nil {
		c.Name = *p.Name
	}
	if p.Email != nil {
		c.Email = *p.Email
	}
	if p.Phone != nil {
		c.Phone = *p.Phone
	}
	if p.Description != nil {
		c.Description = *p.Description
	}
	m.customers[c.ID] = c
	return c, nil
}

func (m *MockStripeClient) GetCustomer(id string) (*stripego.Customer, error) {
	if c, ok := m.customers[id]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("no such customer: %s", id)
}

func (m *MockStripeClient) SearchCustomers(p *stripego.CustomerSearchParams) ([]*stripego.Customer, error) {
	out := make([]*stripego.Customer, 0, len(m.customers))
	for _, c := range m.customers {
		out = append(out, c)
	}
	return out, nil
}

func (m *MockStripeClient) UpdateCustomer(id string, p *stripego.CustomerParams) (*stripego.Customer, error) {
	c, err := m.GetCustomer(id)
	if err != nil {
		return nil, err
	}
	if p.Name != nil {
		c.Name = *p.Name
	}
	if p.Email != nil {
		c.Email = *p.Email
	}
	if p.Phone != nil {
		c.Phone = *p.Phone
	}
	if p.Description != nil {
		c.Description = *p.Description
	}
	return c, nil
}

func (m *MockStripeClient) DeleteCustomer(id string) (*stripego.Customer, error) {
	c, err := m.GetCustomer(id)
	if err != nil {
		return nil, err
	}
	delete(m.customers, id)
	return c, nil
}

// Checkout sessions

func (m *MockStripeClient) CreateSession(p *stripego.CheckoutSessionParams) (*stripego.CheckoutSession, error) {
	id := m.nextID("cs")
	sess := &stripego.CheckoutSession{
		ID:      id,
		URL:     "https://checkout.stripe.com/mock/" + id,
		Mode:    stripego.CheckoutSessionMode("subscription"),
		Status:  stripego.CheckoutSessionStatus("open"),
		Created: time.Now().Unix(),
	}
	m.sessions[sess.ID] = sess
	return sess, nil
}

func (m *MockStripeClient) GetSession(id string) (*stripego.CheckoutSession, error) {
	if s, ok := m.sessions[id]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("no such session: %s", id)
}

func (m *MockStripeClient) ExpireSession(id string) (*stripego.CheckoutSession, error) {
	sess, err := m.GetSession(id)
	if err != nil {
		return nil, err
	}
	sess.Status = stripego.CheckoutSessionStatus("expired")
	return sess, nil
}

func (m *MockStripeClient) ListSessions(_ *stripego.CheckoutSessionListParams) ([]*stripego.CheckoutSession, error) {
	out := make([]*stripego.CheckoutSession, 0, len(m.sessions))
	for _, s := range m.sessions {
		out = append(out, s)
	}
	return out, nil
}

// Subscriptions

func (m *MockStripeClient) GetSubscription(id string) (*stripego.Subscription, error) {
	if s, ok := m.subscriptions[id]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("no such subscription: %s", id)
}

func (m *MockStripeClient) ListSubscriptions(p *stripego.SubscriptionListParams) ([]*stripego.Subscription, error) {
	out := make([]*stripego.Subscription, 0, len(m.subscriptions))
	for _, s := range m.subscriptions {
		if p.Customer != nil && (s.Customer == nil || s.Customer.ID != *p.Customer) {
			continue
		}
		out = append(out, s)
	}
	return out, nil
}
