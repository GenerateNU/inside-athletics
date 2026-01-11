package health

import (
	"context"
	"fmt"
	models "inside-athletics/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"gorm.io/gorm"
)

type HealthDB struct {
	db *gorm.DB
}


func (h *HealthDB) Ping() error {
	return db.Ping()
}
