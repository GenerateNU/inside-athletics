package server

import (
	"context"
	"encoding/json"
	"inside-athletics/internal/handlers/college"
	"inside-athletics/internal/handlers/collegefollow"
	"inside-athletics/internal/handlers/comment"
	"inside-athletics/internal/handlers/comment_like"
	"inside-athletics/internal/handlers/content"
	"inside-athletics/internal/handlers/health"
	"inside-athletics/internal/handlers/media"
	"inside-athletics/internal/handlers/permission"
	"inside-athletics/internal/handlers/post"
	"inside-athletics/internal/handlers/post_like"
	premiumpost "inside-athletics/internal/handlers/premium_post"
	"inside-athletics/internal/handlers/role"
	"inside-athletics/internal/handlers/sport"
	"inside-athletics/internal/handlers/sportfollow"
	"inside-athletics/internal/handlers/stripe"
	"inside-athletics/internal/handlers/survey"
	"inside-athletics/internal/handlers/tag"
	"inside-athletics/internal/handlers/tagfollow"
	"inside-athletics/internal/handlers/tagpost"
	"inside-athletics/internal/handlers/user"
	"inside-athletics/internal/handlers/utility"
	"inside-athletics/internal/s3"
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

// CreateApp initializes the Fiber app and returns the assembled App (server + Huma API).
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

// CreateRoutes registers all route groups on the given Huma API.
func CreateRoutes(db *gorm.DB, api huma.API) {
	api.UseMiddleware(PermissionHumaMiddleware(api, db))
	routeGroups := [...]RouteFN{survey.Route, media.Route, health.Route, sport.Route, role.Route, permission.Route, collegefollow.Route, tagfollow.Route, sportfollow.Route, tagpost.Route, comment.Route, comment_like.Route, post_like.Route, stripe.Route, comment.Route}
	for _, fn := range routeGroups {
		fn(api, db)
	}

	utility.Route(api, db)

	var s3Svc *s3.Service
	if s3Cfg, ok := s3.LoadConfigFromEnv(); ok {
		if client, err := s3.NewClient(context.Background(), s3Cfg); err == nil {
			s3Svc = s3.NewService(client, s3Cfg)
		}
	}

	college.Route(api, db, s3Svc)
	user.Route(api, db, s3Svc)
	post.Route(api, db, s3Svc)
	tag.Route(api, db, s3Svc)
	content.Route(api, db, s3Svc)
	premiumpost.Route(api, db, s3Svc)
}

// setupApp initializes the Fiber app with middleware and returns the configured instance.
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
