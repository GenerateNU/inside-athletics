package stripe

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterWebhookRoute registers the Stripe webhook on the raw Fiber app (bypassing Huma/auth).
func RegisterWebhookRoute(router *fiber.App, db *gorm.DB) {
	svc := NewStripeService(db)
	router.Post("/api/v1/stripe/webhook", svc.HandleWebhook)
}

func RouteWithClient(api huma.API, db *gorm.DB, client StripeClient) {
	registerRoutes(api, NewStripeServiceWithClient(db, client))
}

func Route(api huma.API, db *gorm.DB) {
	registerRoutes(api, NewStripeService(db))
}

func registerRoutes(api huma.API, svc *StripeService) {
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_product")
		huma.Post(grp, "/", svc.CreateStripeProduct)
		huma.Get(grp, "/{id}", svc.GetStripeProductByID)
		huma.Patch(grp, "/{id}", svc.UpdateStripeProduct)
		huma.Delete(grp, "/{id}", svc.ArchiveStripeProduct)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_price")
		huma.Post(grp, "/", svc.CreateStripePrice)
		huma.Get(grp, "/{id}", svc.GetStripePriceByID)
		huma.Patch(grp, "/{id}", svc.UpdateStripePrice)
		huma.Delete(grp, "/{id}", svc.ArchiveStripePrice)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_products")
		huma.Get(grp, "/", svc.GetAllStripeProducts)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_prices")
		huma.Get(grp, "/{id}", svc.GetAllStripePrices)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_customers")
		huma.Post(grp, "/", svc.RegisterStripeCustomer)
		huma.Get(grp, "/{id}", svc.GetStripeCustomer)
		huma.Get(grp, "/email/{email}", svc.GetStripeCustomerByEmail)
		huma.Patch(grp, "/{id}", svc.UpdateStripeCustomer)
		huma.Delete(grp, "/{id}", svc.DeleteStripeCustomer)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/checkout/sessions")
		huma.Post(grp, "/", svc.CreateStripeCheckoutSession)
		huma.Get(grp, "/", svc.GetAllStripeSessions)
		huma.Get(grp, "/{id}", svc.GetStripeCheckoutSessionByID)
		huma.Delete(grp, "/{id}", svc.DeleteStripeCheckoutSession)
	}
}
