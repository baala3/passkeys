package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

const webauthnSessionDuration = 30 * 24 * time.Hour

type WebAuthnSession struct {
	RedisClient *redis.Client
}

func (session *WebAuthnSession) Get(ctx echo.Context, sessionName string) (string,*webauthn.SessionData, error) {
	cookie, err := ctx.Cookie(sessionName)

	if err != nil {
		return "", nil, fmt.Errorf("failed to get session data: %v", err)
	}

	id := cookie.Value

	bytes, err := session.RedisClient.Get(ctx.Request().Context(), id).Bytes()

	if err != nil {
		return "", nil, fmt.Errorf("failed to get session data: %v", err)
	}

	// Unmarshal session data from JSON
	var data *webauthn.SessionData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return "", nil, fmt.Errorf("failed to unmarshal session data: %v", err)
	}
	return id, data, nil
}

func (session *WebAuthnSession) Create(ctx echo.Context, sessionName string, data *webauthn.SessionData) ( error) {
	// Marshal session data to JSON
	bytes, err := json.Marshal(data)
	if err != nil {
		return  fmt.Errorf("failed to encode session data: %v", err)	
	}

	id := uuid.New().String()

	if err := session.RedisClient.Set(ctx.Request().Context(), id, bytes, webauthnSessionDuration).Err(); err != nil {
		return fmt.Errorf("failed to save session data: %v", err)
	}

	ctx.SetCookie(&http.Cookie{
		Name: sessionName,
		Value: id,
		Path: "/",
	})

	return nil
}

func (session *WebAuthnSession) Delete(ctx echo.Context, id string) error {
	if err := session.RedisClient.Del(ctx.Request().Context(), id).Err(); err != nil {
		return fmt.Errorf("failed to delete session data: %v", err)
	}
	return nil
}