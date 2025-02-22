package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func ConditionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		context := c.QueryParam("context")
		switch context {
		case "signup":
			return next(c)
		default:
			return Auth(next)(c)
		}
	}
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		redisClient:= ctx.Get("redisClient").(*redis.Client)

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

func NoAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		redisClient:= ctx.Get("redisClient").(*redis.Client)
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
