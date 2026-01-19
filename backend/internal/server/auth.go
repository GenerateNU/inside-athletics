package server

import (
	"context"
	"strings"
	"net/http"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/danielgtaylor/huma/v2"
)


/**
Middleware used for verifying JWT in Authorization header
and injecting userID into context for use across all API
endpoints
*/
func AuthMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// The url holding the public key. This can be exposed publicly as it is only for
		// verification purposes and should be used as a URL so we can rotate keys
		jwksURL := "https://zhhniddxrmfqqracjrlc.supabase.co/auth/v1/.well-known/jwks.json"

		// created a new keyfunc.KeyFunc which is used
		// when verifying the token 
		bgContext := context.Background()
		key, err := keyfunc.NewDefaultCtx(bgContext, []string{jwksURL})
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError,
				"Unable to verify JWT", err,
			)
			return
		}

		// extracting the JWT from the authorization
		// header. If it cannot be extracted an Unauthorized
		// error will be thrown
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized,
				"Authorization Header not in Request", err,
			)
			return
		}
		headerComponents := strings.Split(authHeader, "")
		if len(headerComponents) != 2 && headerComponents[0] != "Bearer"{
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized,
				"Bearer not ioncluded in Authorization Header",
			)
			return
		}
		token := headerComponents[1]

		// Parse the jwt to see if the key has been maniuplated with or
		// is not signed with the correct private key
		parsed, parseErr := jwt.Parse(token, key.Keyfunc, jwt.WithValidMethods([]string{"ES256"}))
		if parseErr != nil  || !parsed.Valid {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized,
				"Token is not valid", parseErr,
			)
			return
		}

		// Now that the key is verified we extract out the userID
		// and save this in the context so that we can access it in
		// any method
		claims, parseWorks := parsed.Claims.(jwt.MapClaims)
		if !parseWorks {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError,
				"Unable to extract User ID",
			)
			return
		}

		userID, userFetched := claims["sub"].(string)
		if !userFetched {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError,
				"Unable to extract User ID",
			)
			return
		}
		ctx = huma.WithValue(ctx, "user_id", userID)

		next(ctx)
}
}