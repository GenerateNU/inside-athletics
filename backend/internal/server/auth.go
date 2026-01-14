package server

import (
	"context"
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

	// created a new keyfunc.KeyFunc which is used
	// when verifying the token 
	ctx := context.Background()
	key, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return fiber.ErrInternalServerError;
	}

	// extracting the JWT from the authorization
	// header. If it cannot be extracted an Unauthorized
	// error will be thrown
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return fiber.ErrUnauthorized;
	}
	headerComponents := strings.Split(authHeader, "")
	if len(headerComponents) != 3 && headerComponents[1] != "Bearer"{
		return fiber.ErrUnauthorized;
	}
	token := headerComponents[2]

	// Parse the jwt to see if the key has been maniuplated with or
	// is not signed with the correct private key
	parsed, parseErr := jwt.Parse(token, key.Keyfunc, jwt.WithValidMethods([]string{"ES256"}))
	if parseErr != nil  || !parsed.Valid {
		return fiber.ErrUnauthorized
	}

	// Now that the key is verified we extract out the userID
	// and save this in the context so that we can access it in
	// any method
	claims, parseWorks := parsed.Claims.(jwt.MapClaims)
	if !parseWorks {
		return fiber.ErrInternalServerError
	}

	userID, userFetched := claims["sub"].(string)
	if !userFetched {
		return fiber.ErrInternalServerError
	}
	c.Set("userID", userID);

	return c.Next()
}