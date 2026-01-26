package server

import (
    "context"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/MicahParks/keyfunc/v3"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

// contextKey is a private type to avoid collisions in context keys.
type contextKey string

/**
Middleware used for verifying JWT in Authorization header
and injecting userID into context for use across all API
endpoints
*/
func AuthMiddleware(c *fiber.Ctx) error {
    env := os.Getenv("APP_ENV") 
    var jwksURL string
    if env == "production" {
		jwksURL = "https://zhhniddxrmfqqracjrlc.supabase.co/auth/v1/.well-known/jwks.json"
	} else {
		jwksURL = "http://localhost:54321/auth/v1/.well-known/jwks.json"
	}
    
    bgContext := context.Background()
    
	key, err := keyfunc.NewDefaultCtx(bgContext, []string{jwksURL})
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Unable to Verify JWT",
    }	)
	}


    authHeader := c.Get("Authorization")
    log.Print(authHeader)
    if authHeader == "" {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Authorization header not in request",
        })
    }
    
    headerComponents := strings.Split(authHeader, " ")
    log.Print(headerComponents)
    if len(headerComponents) != 2 || headerComponents[0] != "Bearer" {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
            "error": "Bearer not included in Authorization header",
        })
    }
    
    token := headerComponents[1]
    
    parsed, parseErr := jwt.Parse(token, key.Keyfunc, jwt.WithValidMethods([]string{"ES256"}))
    if parseErr != nil || !parsed.Valid {
        log.Print(parseErr)
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
    c.SetUserContext(context.WithValue(c.UserContext(), contextKey("user_id"), userID))
    return c.Next()
}
