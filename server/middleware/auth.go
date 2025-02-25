package middleware

import (
	"context"
	"net/http"

	"github.com/baala3/passkeys/pkg"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// ConditionalAuth allows passkey-signup without authentication, otherwise applies Auth middleware
func ConditionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		context := c.QueryParam("context")
		switch context {
		case "signup", "signin":
			return next(c)
		default:
			return Auth(next)(c)
		}
	}
}

// Auth middleware ensures user is authenticated
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		redisClient := pkg.GetRedisClient()
		if redisClient == nil {
			return ctx.Redirect(http.StatusFound, "/")
		}

		cookie, err := ctx.Cookie("auth")
		if err != nil {
			return ctx.Redirect(http.StatusFound, "/")
		}

		sessionID := cookie.Value
		userID, err := redisClient.Get(context.Background(), sessionID).Result()
		if err == redis.Nil {
			return ctx.Redirect(http.StatusFound, "/")
		}

		if err != nil {
			return ctx.Redirect(http.StatusFound, "/")
		}
		// // Store userID in context for later use
		ctx.Set("userID", userID)
		return next(ctx)
	}
}

// NoAuth ensures the user is NOT authenticated
func NoAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		redisClient := pkg.GetRedisClient()
		if redisClient == nil {
			return next(ctx)
		}
		cookie, err := ctx.Cookie("auth")
		if err != nil {
			return next(ctx)
		}
		sessionID := cookie.Value
		userID, err := redisClient.Get(context.Background(), sessionID).Result()
		if err == redis.Nil || err != nil || userID == "" {
			return next(ctx)
		}
		return ctx.Redirect(http.StatusFound, "/home")
	}
}
