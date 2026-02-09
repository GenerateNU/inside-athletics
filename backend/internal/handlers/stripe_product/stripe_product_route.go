package stripe_product

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	stripeService := NewStripeProductService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_product")
		huma.Post(grp, "/", stripeService.CreateStripeProduct)                // Create product (subscription plan)
		huma.Get(grp, "/{id}", stripeService.GetStripeProductByID)            // Read product by ID
		huma.Patch(grp, "/{id}", stripeService.UpdateStripeProduct)           // Update product
		huma.Delete(grp, "/{id}", stripeService.DeleteStripeProduct)          // Delete product
	}
	{
		grp := huma.NewGroup(api, "/api/v1/stripe_products")
		huma.Get(grp, "/", stripeService.GetAllStripeProducts)                // Create all products (subscription plan)
	}
}
