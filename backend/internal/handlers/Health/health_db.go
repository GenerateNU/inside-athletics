package health

import (
	"context"
	"fmt"
	models "inside-athletics/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthDB struct {
	conn *pgxpool.Pool
}

func (h *HealthDB) GetFromDB(name string) (*models.HealthModel, error) {
	sql := fmt.Sprintf("select * from \"Test Table\" where name = '%s'", name)

	rows, err := h.conn.Query(context.Background(), sql)

	if err != nil {
		return &models.HealthModel{}, huma.Error404NotFound(fmt.Sprintf("Could not fetch value with id %s, %s", name, err.Error()))
	}

	healthModels, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.HealthModel])

	if err != nil {
		return &models.HealthModel{}, huma.Error400BadRequest("Unable to collect row into HealthModel object, the given id is not stored in the db")
	}

	return &healthModels, nil
}
