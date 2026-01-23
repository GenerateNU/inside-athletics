package sport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SportHandler struct {
	db *SportDB
}

func NewSportHandler(db *gorm.DB) *SportHandler {
	return &SportHandler{
		db: &SportDB{db: db},
	}
}

type CreateSportRequest struct {
	Name       string `json:"name" binding:"required"`
	Popularity int32  `json:"popularity" binding:"required"`
}

func (h *SportHandler) CreateSport(c *gin.Context) {
	var req CreateSportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sport, err := h.db.CreateSport(req.Name, req.Popularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sport)
}

func (h *SportHandler) GetSportById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	sport, err := h.db.GetSportById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sport)
}

func (h *SportHandler) GetAllSports(c *gin.Context) {
	sports, err := h.db.GetAllSports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sports)
}

func (h *SportHandler) UpdateSport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var req CreateSportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sport, err := h.db.UpdateSport(id, req.Name, req.Popularity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sport)
}

func (h *SportHandler) DeleteSport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	sport, err := h.db.DeleteSport(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sport)
}
