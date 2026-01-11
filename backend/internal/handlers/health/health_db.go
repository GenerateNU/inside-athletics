package health

import (
	"gorm.io/gorm"
)

type HealthDB struct {
	db *gorm.DB
}


func (h *HealthDB) Ping() error {
	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
