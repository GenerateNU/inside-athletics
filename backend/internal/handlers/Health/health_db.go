package health

import (
	"context"
	"fmt"
	models "inside-athletics/internal/models"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthDB struct {
	conn *pgxpool.Pool
}

func (h *HealthDB) GetFromDB(id string) (*models.HealthModel, error) {
	sql := fmt.Sprintf("select * from \"Test Table\" where id = %s", id)

	rows, err := h.conn.Query(context.Background(), sql)

	if err != nil {
		return &models.HealthModel{}, huma.NewError(http.StatusBadRequest, fmt.Sprintf("Could not fetch value with id %s", id))
	}

	healthModels, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.HealthModel])

	if err != nil {
		return &models.HealthModel{}, huma.NewError(http.StatusBadRequest, "Unable to collect row into HealthModel object, the given id is not stored in the db")
	}

	return &healthModels, nil
}
