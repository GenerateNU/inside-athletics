package stripe

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	stripeService := NewStripeService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_product")
		huma.Post(grp, "/", stripeService.CreateStripeProduct)        // Create product (subscription plan)
		huma.Get(grp, "/{id}", stripeService.GetStripeProductByID)    // Read product by ID
		huma.Patch(grp, "/{id}", stripeService.UpdateStripeProduct)   // Update product
		huma.Delete(grp, "/{id}", stripeService.ArchiveStripeProduct) // Delete product
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_price")
		huma.Post(grp, "/", stripeService.CreateStripePrice)        // Create price for specified product
		huma.Get(grp, "/{id}", stripeService.GetStripePriceByID)    // Read price by ID
		huma.Patch(grp, "/{id}", stripeService.UpdateStripePrice)   // Update price
		huma.Delete(grp, "/{id}", stripeService.ArchiveStripePrice) // Delete price
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_products")
		huma.Get(grp, "/", stripeService.GetAllStripeProducts) // Get all products (subscription plan)
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_prices")
		huma.Get(grp, "/{id}", stripeService.GetAllStripePrices) // Get all prices from specific product
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_customers")
		huma.Post(grp, "/", stripeService.RegisterStripeCustomer)     // Register a new customer
		huma.Get(grp, "/{id}", stripeService.GetStripeCustomer)       // Get a customer
		huma.Patch(grp, "/{id}", stripeService.UpdateStripeCustomer)  // Update a customer
		huma.Delete(grp, "/{id}", stripeService.DeleteStripeCustomer) // Delete a customer
	}
}
