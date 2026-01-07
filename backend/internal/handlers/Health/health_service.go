package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthService struct {
	healthDB *HealthDB
}

func (h *HealthService) CheckHealth(c *gin.Context) {
	c.JSON(200, gin.H{"message": "AYYYY"})
}

func (h *HealthService) GetHealth(c *gin.Context) {
	id := c.Param("id")

	healthModel, err := h.healthDB.GetFromDB(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(200, healthModel)
}
