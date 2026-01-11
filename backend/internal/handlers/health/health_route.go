package health

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func Route(api huma.API, db *gorm.DB) {
	var healthDB *HealthDB = &HealthDB{db}              // create object storing all database level functions for health
	var healthService *HealthService = &HealthService{healthDB} // create object with health functionality
	{
		grp := huma.NewGroup(api, "/api/v1/health")
		huma.Get(grp, "/", healthService.Health)
		huma.Get(grp, "/healthcheck", healthService.CheckHealth)
	}
}
