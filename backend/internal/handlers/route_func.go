package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouteFN func(app *gin.Engine, connection *pgxpool.Pool)
