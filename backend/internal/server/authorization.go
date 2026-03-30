package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"inside-athletics/internal/models"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	errInvalidResourceID   = errors.New("invalid resource ID")
	errUnsupportedResource = errors.New("unsupported resource for ownership check")
)

var resourceByPathPrefix = map[string]string{
	"user/tag":     "tagfollow",
	"user/sport":   "sportfollow",
	"user/college": "collegefollow",
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

		parsedUserID, err := uuid.Parse(userID)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Invalid user ID")
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

		if (action == models.PermissionUpdate || action == models.PermissionDelete) &&
			(resource == "post" || resource == "comment" || resource == "tagfollow" || resource == "premiumpost" || resource == "sportfollow" || resource == "collegefollow") &&
			ctx.Param("id") != "" {
			owned, err := IsOwnerOfResource(db, parsedUserID, ctx.Param("id"), resource)
			if err != nil {
				status := http.StatusInternalServerError
				if errors.Is(err, errInvalidResourceID) || errors.Is(err, errUnsupportedResource) {
					status = http.StatusBadRequest
				}
				_ = huma.WriteErr(api, ctx, status, err.Error())
				return
			}
			if owned {
				if action == models.PermissionDelete {
					action = models.PermissionDeleteOwn
				} else {
					action = models.PermissionUpdateOwn
				}
			}
		}

		allowed, status, msg := authorizeByPermission(db, parsedUserID, action, resource)
		if !allowed {
			_ = huma.WriteErr(api, ctx, status, msg)
			return
		}

		next(ctx)
	}
}

func authorizeByPermission(db *gorm.DB, userID uuid.UUID, action models.PermissionAction, resource string) (bool, int, string) {
	authDB := NewAuthorizationDB(db)

	if err := authDB.UserExists(userID); err != nil {
		return false, http.StatusUnauthorized, "User not found"
	}

	hasPermission, err := authDB.UserHasPermission(userID, action, resource)
	if err != nil {
		return false, http.StatusInternalServerError, "Unable to check permissions"
	}

	if !hasPermission {
		return false, http.StatusForbidden, "Insufficient permissions"
	}

	return true, 0, ""
}

func IsOwnerOfResource(db *gorm.DB, userID uuid.UUID, resourceID, resource string) (bool, error) {
	parsedResourceID, err := uuid.Parse(resourceID)
	if err != nil {
		return false, errInvalidResourceID
	}

	var count int64
	switch resource {
	case "post":
		err = db.Table("posts").
			Where("id = ? AND author_id = ?", parsedResourceID, userID).
			Count(&count).Error
	case "comment":
		err = db.Table("comments").
			Where("id = ? AND user_id = ?", parsedResourceID, userID).
			Count(&count).Error
	case "tagfollow":
		err = db.Table("tag_follows").
			Where("id = ? AND user_id = ?", parsedResourceID, userID).
			Count(&count).Error
	case "sportfollow":
		err = db.Table("sport_follows").
			Where("id = ? AND user_id = ?", parsedResourceID, userID).
			Count(&count).Error
	case "collegefollow":
		err = db.Table("college_follows").
			Where("id = ? AND user_id = ?", parsedResourceID, userID).
			Count(&count).Error
	case "premiumpost":
		err = db.Table("premium_posts").
			Where("id = ? AND author_id = ?", parsedResourceID, userID).
			Count(&count).Error
	default:
		return false, errUnsupportedResource
	}

	if err != nil {
		return false, fmt.Errorf("unable to check ownership: %w", err)
	}
	return count > 0, nil
}

func IsOwnerOfPostOrComment(db *gorm.DB, userID uuid.UUID, resourceID, resource string) (bool, error) {
	return IsOwnerOfResource(db, userID, resourceID, resource)
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
	"user":          "user",
	"users":         "user",
	"post":          "post",
	"posts":         "post",
	"sport":         "sport",
	"sports":        "sport",
	"college":       "college",
	"colleges":      "college",
	"role":          "role",
	"roles":         "role",
	"permission":    "permission",
	"permissions":   "permission",
	"premium_post":  "premium_post",
	"premium_posts": "premium_post",
}

func resolveResourceFromPath(path string) string {
	path = strings.Trim(strings.TrimPrefix(path, "/api/v1/"), "/")
	for prefix, resource := range resourceByPathPrefix {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return resource
		}
	}

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
