package pkg

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/redis/go-redis/v9"
)

const sessionDuration = 30 * 24 * time.Hour

type UserSession struct {
	RedisClient *redis.Client
}

func (ss *UserSession) Create(ctx echo.Context, userID uuid.UUID) error {
	sessionID := random.String(20)

	if err := ss.RedisClient.Set(ctx.Request().Context(), sessionID, userID.String(), sessionDuration).Err(); err != nil {
		return fmt.Errorf("failed to save session data: %v", err)
	}

	ctx.SetCookie(&http.Cookie{
		Name: "auth",
		Value: sessionID,
		Path: "/",
	})
	return nil
}

func (ss *UserSession) Delete(ctx echo.Context, sessionID string) {
	_ = ss.RedisClient.Del(ctx.Request().Context(), sessionID).Err()
}
