package health

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Route(app *gin.Engine, connection *pgxpool.Pool) {
	var healthDB *HealthDB = &HealthDB{connection}              // create object storing all database level functions for health
	var healthService *HealthService = &HealthService{healthDB} // create object with health functionality

	{
		group := app.Group("/api/v1/health")
		group.GET("/", healthService.CheckHealth)
		group.GET("/:id", healthService.GetHealth)
	}
}
