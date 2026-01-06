package health

import (
	"context"
	"errors"
	"fmt"
	health "inside-athletics/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthDB struct {
	conn *pgxpool.Pool
}

func (h *HealthDB) GetFromDB(id string) (health.HealthModel, error) {
	sql := fmt.Sprintf("select * from \"Test Table\" where id = %s", id)

	rows, err := h.conn.Query(context.Background(), sql)

	if err != nil {
		return health.HealthModel{}, fmt.Errorf("Could not fetch value with id %s", id)
	}

	healthModels, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[health.HealthModel])

	if err != nil {
		return health.HealthModel{}, errors.New("Unable to collect row into HealthModel object, the given id is not stored in the db")
	}

	return healthModels, nil
}
