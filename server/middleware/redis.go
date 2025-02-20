package middleware

import (
	"net/http"

	"github.com/baala3/passkeys/pkg"
	"github.com/labstack/echo/v4"
)

func InjectRedis(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		redisClient := pkg.GetRedisClient()
		if redisClient == nil {
			return c.Redirect(http.StatusFound, "/")
		}
		c.Set("redisClient", redisClient)
		return next(c)
	}
}

