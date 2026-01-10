package health

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Route(api huma.API, connection *pgxpool.Pool) {
	var healthDB *HealthDB = &HealthDB{connection}              // create object storing all database level functions for health
	var healthService *HealthService = &HealthService{healthDB} // create object with health functionality

	huma.Get(api, "/health", healthService.CheckHealth)
}
