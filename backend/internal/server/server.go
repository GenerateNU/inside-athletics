package server

import (
	"encoding/json"
	"inside-athletics/internal/handlers/college"
	"inside-athletics/internal/handlers/health"
	"inside-athletics/internal/handlers/sport"
	"inside-athletics/internal/handlers/stripe"
	"inside-athletics/internal/handlers/user"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/skip"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"

	"gorm.io/gorm"
)

type App struct {
	Server *fiber.App
	Api    huma.API
}

type RouteFN func(api huma.API, db *gorm.DB)

func CreateApp(db *gorm.DB) *App {

	router := setupApp()
	config := huma.DefaultConfig("Inside Athletics API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"Authorization": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	config.Security = []map[string][]string{
		{
			"Authorization": {},
		},
	}

	var api = humafiber.New(router, config)
	CreateRoutes(db, api)
	return &App{
		Server: router,
		Api:    api,
	}
}

func CreateRoutes(db *gorm.DB, api huma.API) {
	// Create all the routing groups:
	routeGroups := [...]RouteFN{health.Route, user.Route, sport.Route, college.Route, stripe.Route}
	for _, fn := range routeGroups {
		fn(api, db)
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

	app.Use(skip.New(AuthMiddleware, func(ctx *fiber.Ctx) bool {
		return strings.HasPrefix(ctx.Path(), "/docs") || strings.HasPrefix(ctx.Path(), "/openapi.yaml") || ctx.Path() == "/"
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
