package server

import (
	"context"
	"net/http"
	"strings"

	"inside-athletics/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

		action, resource := resolveResourceAndAction(ctx.Method(), path, userID, ctx.Param("id"))
		if action == "" || resource == "" {
			next(ctx)
			return
		}

		allowed, status, msg := authorizeByPermission(db, userID, action, resource)
		if !allowed {
			_ = huma.WriteErr(api, ctx, status, msg)
			return
		}

		next(ctx)
	}
}

func authorizeByPermission(db *gorm.DB, userID string, action models.PermissionAction, resource string) (bool, int, string) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return false, http.StatusUnauthorized, "Invalid user ID"
	}

	var user models.User
	if err := db.Select("id").First(&user, "id = ?", parsedUserID).Error; err != nil {
		return false, http.StatusUnauthorized, "User not found"
	}

	var count int64
	if err := db.Table("user_roles").
		Joins("JOIN role_permissions rp ON rp.role_id = user_roles.role_id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("user_roles.user_id = ? AND p.action = ? AND p.resource = ?", user.ID, action, resource).
		Count(&count).Error; err != nil {
		return false, http.StatusInternalServerError, "Unable to check permissions"
	}

	if count == 0 {
		return false, http.StatusForbidden, "Insufficient permissions"
	}

	return true, 0, ""
}

func resolveResourceAndAction(method, path, userID, pathID string) (models.PermissionAction, string) {
	action := models.PermissionAction("")
	switch method {
	case http.MethodPost:
		action = models.PermissionCreate
	case http.MethodPut, http.MethodPatch:
		action = models.PermissionUpdate
	case http.MethodDelete:
		action = models.PermissionDelete
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
				action = models.PermissionDeleteOwn
			} else {
				action = models.PermissionUpdateOwn
			}
		}
	}

	return action, resource
}

// change this to map so that you can just change the map if you need to update
func resolveResourceFromPath(path string) string {
	path = strings.TrimPrefix(path, "/api/v1/")
	segment := strings.SplitN(path, "/", 2)[0]
	switch segment {
	case "user", "users":
		return "user"
	case "post", "posts":
		return "post"
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
