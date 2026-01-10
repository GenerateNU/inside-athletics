package server

import (
	"encoding/json"
	"inside-athletics/internal/handlers"
	health "inside-athletics/internal/handlers/Health"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
)

type App struct {
	Server *fiber.App
	Api    huma.API
}

func CreateApp(connection *pgxpool.Pool) *App {

	router := setupApp()
	var api huma.API = humafiber.New(router, huma.DefaultConfig("Inside Athletics API", "1.0.0"))

	// Create all the routing groups:
	routeGroups := [...]handlers.RouteFN{health.Route}
	for _, fn := range routeGroups {
		fn(api, connection)
	}

	return &App{
		Server: router,
		Api:    api,
	}
}

// Initialize Fiber app with middlewares / configs
func setupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}:${port} ${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(favicon.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	return app
}
