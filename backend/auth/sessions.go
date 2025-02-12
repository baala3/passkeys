package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/redis/go-redis/v9"
)

var sessionStore *redis.Client
const duration time.Duration = 5 * time.Minute

func init() {
	sessionStore = redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
		Password: "",
		DB: 0,
	})
}

func GetSession(ctx context.Context, username string) (*webauthn.SessionData, error) {
	bytes, err := sessionStore.Get(ctx, username).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get session data: %v", err)
	}
	// Unmarshal session data from JSON
	var data *webauthn.SessionData
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %v", err)
	}
	return data, nil
}

func CreateSession(ctx context.Context, username string, data *webauthn.SessionData) error {
	// Marshal session data to JSON
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to encode session data: %v", err)	
	}

	if err := sessionStore.Set(ctx, username, bytes, duration).Err(); err != nil {
		return fmt.Errorf("failed to save session data: %v", err)
	}

	return nil
}
