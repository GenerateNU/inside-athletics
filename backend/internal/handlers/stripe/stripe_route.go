package stripe

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	stripeService := NewStripeService(db)
	{
		grp := huma.NewGroup(api, "/api/v1/stripe")
		huma.Post(grp, "/", stripeService.CreateProduct)                // Create product (subscription plan)
		huma.Get(grp, "/{id}", stripeService.GetProdcutByID)            // Read product by ID
		huma.Patch(grp, "/{id}", stripeService.UpdateProduct)           // Update product
		huma.Delete(grp, "/{id}", stripeService.DeleteProduct)          // Delete product
	}
}
