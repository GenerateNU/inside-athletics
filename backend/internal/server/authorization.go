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
		if resource == "user" && action == models.PermissionCreate {
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

	authDB := NewAuthorizationDB(db)

	if err := authDB.UserExists(parsedUserID); err != nil {
		return false, http.StatusUnauthorized, "User not found"
	}

	hasPermission, err := authDB.UserHasPermission(parsedUserID, action, resource)
	if err != nil {
		return false, http.StatusInternalServerError, "Unable to check permissions"
	}

	if !hasPermission {
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

	// Special-case: if the user is acting on their own resource, use the *_own permissions.
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
var resourceByPathSegment = map[string]string{
	"user":        "user",
	"users":       "user",
	"post":        "post",
	"posts":       "post",
	"sport":       "sport",
	"sports":      "sport",
	"college":     "college",
	"colleges":    "college",
	"role":        "role",
	"roles":       "role",
	"permission":  "permission",
	"permissions": "permission",
}

func resolveResourceFromPath(path string) string {
	path = strings.TrimPrefix(path, "/api/v1/")
	segment := strings.SplitN(path, "/", 2)[0]
	return resourceByPathSegment[segment]
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
