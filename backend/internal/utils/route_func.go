package utils

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouteFN func(api huma.API, connection *pgxpool.Pool)
