package server

import (
	"context"
	"net/http"
	"strings"

	"inside-athletics/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AdminOnlyMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return c.Next()
		}

		userIDValue := c.Locals("user_id")
		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var user models.User
		if err := db.Preload("Role").First(&user, "id = ?", parsedUserID).Error; err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		if user.Role.Name != models.RoleAdmin {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Admin privileges required",
			})
		}

		return c.Next()
	}
}

func PermissionMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch c.Method() {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return c.Next()
		}

		userIDValue := c.Locals("user_id")
		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		allowed, status, msg := authorizeByPermission(db, userID, c.Method(), c.Path(), c.Params("id"))
		if !allowed {
			return c.Status(status).JSON(fiber.Map{
				"error": msg,
			})
		}

		return c.Next()
	}
}

func PermissionHumaMiddleware(api huma.API, db *gorm.DB) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		switch ctx.Method() {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next(ctx)
			return
		}

		path := ctx.URL().Path
		if strings.HasPrefix(path, "/docs") || strings.HasPrefix(path, "/openapi.yaml") || path == "/" {
			next(ctx)
			return
		}

		userID, ok := getUserIDFromContext(ctx.Context())
		if !ok || userID == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "User not authenticated")
			return
		}

		allowed, status, msg := authorizeByPermission(db, userID, ctx.Method(), path, ctx.Param("id"))
		if !allowed {
			_ = huma.WriteErr(api, ctx, status, msg)
			return
		}

		next(ctx)
	}
}

func authorizeByPermission(db *gorm.DB, userID, method, path, pathID string) (bool, int, string) {
	resource, action := resolveResourceAndAction(method, path, userID, pathID)
	if resource == "" || action == "" {
		return true, 0, ""
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return false, http.StatusUnauthorized, "Invalid user ID"
	}

	var perm models.Permission
	if err := db.Select("id").Where("action = ? AND resource = ?", action, resource).First(&perm).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, 0, ""
		}
		return false, http.StatusInternalServerError, "Unable to check permissions"
	}

	var user models.User
	if err := db.Select("role_id").First(&user, "id = ?", parsedUserID).Error; err != nil {
		return false, http.StatusUnauthorized, "User not found"
	}

	var count int64
	if err := db.Table("role_permissions").
		Joins("JOIN permissions p ON p.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND p.action = ? AND p.resource = ?", user.RoleID, action, resource).
		Count(&count).Error; err != nil {
		return false, http.StatusInternalServerError, "Unable to check permissions"
	}

	if count == 0 {
		return false, http.StatusForbidden, "Insufficient permissions"
	}

	return true, 0, ""
}

func resolveResourceAndAction(method, path, userID, pathID string) (string, string) {

	action := ""
	switch method {
	case http.MethodPost:
		action = string(models.PermissionCreate)
	case http.MethodPut, http.MethodPatch:
		action = string(models.PermissionUpdate)
	case http.MethodDelete:
		action = string(models.PermissionDelete)
	default:
		return "", ""
	}

	resource := resolveResourceFromPath(path)
	if resource == "" {
		return "", ""
	}

	if resource == "user" && (method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete) {
		if userID != "" && pathID != "" && pathID == userID {
			if method == http.MethodDelete {
				action = string(models.PermissionDeleteOwn)
			} else {
				action = string(models.PermissionUpdateOwn)
			}
		}
	}

	return resource, action
}

func resolveResourceFromPath(path string) string {
	if strings.HasPrefix(path, "/api/v1/") {
		path = strings.TrimPrefix(path, "/api/v1/")
	}
	segment := strings.SplitN(path, "/", 2)[0]
	switch segment {
	case "user", "users":
		return "user"
	case "sport", "sports":
		return "sport"
	case "college", "colleges":
		return "college"
	case "role", "roles":
		return "role"
	case "permission", "permissions":
		return "permission"
	default:
		return ""
	}
}

func getUserIDFromContext(ctx context.Context) (string, bool) {
	if v, ok := ctx.Value("user_id").(string); ok && v != "" {
		return v, true
	}
	if v, ok := ctx.Value(contextKey("user_id")).(string); ok && v != "" {
		return v, true
	}
	return "", false
}
