package sport

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SportHandler struct {
	db *SportDB
}

func NewSportHandler(db *gorm.DB) *SportHandler {
	return &SportHandler{
		db: NewSportDB(db),
	}
}

// Error response helpers
func (h *SportHandler) respondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func (h *SportHandler) handleDBError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		h.respondWithError(c, http.StatusNotFound, "Sport not found")
		return
	}
	h.respondWithError(c, http.StatusInternalServerError, err.Error())
}

func (h *SportHandler) CreateSport(c *gin.Context) {
	var req CreateSportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	sport, err := h.db.CreateSport(req.Name, req.Popularity)
	if err != nil {
		h.handleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, ToSportResponse(sport))
}

func (h *SportHandler) GetSportByID(c *gin.Context) {
	var params GetSportByIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid sport ID")
		return
	}

	sport, err := h.db.GetSportByID(params.ID)
	if err != nil {
		h.handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, ToSportResponse(sport))
}

func (h *SportHandler) GetAllSports(c *gin.Context) {
	var params GetAllSportsParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	sports, total, err := h.db.GetAllSports(params.Limit, params.Offset)
	if err != nil {
		h.handleDBError(c, err)
		return
	}

	sportResponses := make([]SportResponse, 0, len(sports))
	for i := range sports {
		sportResponses = append(sportResponses, *ToSportResponse(&sports[i]))
	}

	c.JSON(http.StatusOK, GetAllSportsResponse{
		Sports: sportResponses,
		Total:  int(total),
	})
}

func (h *SportHandler) UpdateSport(c *gin.Context) {
	var params GetSportByIDParams
	if err := c.ShouldBindUri(&params); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid sport ID")
		return
	}

	var req UpdateSportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	existingSport, err := h.db.GetSportByID(params.ID)
	if err != nil {
		h.handleDBError(c, err)
		return
	}

	// Apply partial updates
	if req.Name != nil {
		existingSport.Name = *req.Name
	}
	if req.Popularity != nil {
		existingSport.Popularity = req.Popularity
	}

	sport, err := h.db.UpdateSport(existingSport)
	if err != nil {
		h.handleDBError(c, err)
		return
	}

	c.JSON(http.StatusOK, ToSportResponse(sport))
}

func (h *SportHandler) DeleteSport(c *gin.Context) {
	var params DeleteSportRequest
	if err := c.ShouldBindUri(&params); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "Invalid sport ID")
		return
	}

	err := h.db.DeleteSport(params.ID)
	if err != nil {
		h.handleDBError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
