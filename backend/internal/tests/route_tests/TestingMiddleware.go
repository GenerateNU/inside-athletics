package routeTests

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"strings"
)

/*
*
Mock Testing Middleware used for verifying JWT in Authorization header
and injecting userID into context for use across all API
endpoints
*/
func MockAuthMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		// extracting the JWT from the authorization
		// header. If it cannot be extracted an Unauthorized
		// error will be thrown
		authHeader := ctx.Header("Authorization")
		if authHeader == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized,
				"Authorization Header not in Request",
			)
			return
		}
		headerComponents := strings.Split(authHeader, " ")
		fmt.Println("headerComponents:", headerComponents)
		if len(headerComponents) != 2 || headerComponents[0] != "Bearer" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized,
				"Bearer not included in Authorization Header",
			)
			return
		}
		userID := headerComponents[1]
		ctx = huma.WithValue(ctx, "user_id", userID)

		next(ctx)
	}
}
