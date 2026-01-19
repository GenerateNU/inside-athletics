package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

/**
Middleware used for verifying JWT in Authorization header
and injecting userID into context for use across all API
endpoints
*/
func AuthMiddleware(c *fiber.Ctx) error {

	jwksURL := "https://zhhniddxrmfqqracjrlc.supabase.co/auth/v1/.well-known/jwks.json"
    bgContext := context.Background()
    
	key, err := keyfunc.NewDefaultCtx(bgContext, []string{jwksURL})
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unable to Verify JWT",
    }	)
	}


    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Authorization header not in request",
        })
    }
    
    headerComponents := strings.Split(authHeader, " ")
    if len(headerComponents) != 2 || headerComponents[0] != "Bearer" {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Bearer not included in Authorization header",
        })
    }
    
    token := headerComponents[1]
    
    parsed, parseErr := jwt.Parse(token, key.Keyfunc, jwt.WithValidMethods([]string{"ES256"}))
    if parseErr != nil || !parsed.Valid {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Token is not valid",
        })
    }
    
    claims, ok := parsed.Claims.(jwt.MapClaims)
    if !ok {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to extract claims",
        })
    }
    
    userID, ok := claims["sub"].(string)
    if !ok {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": "Unable to extract user ID",
        })
    }
    
    c.Locals("user_id", userID)
    return c.Next()
}