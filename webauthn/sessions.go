package webauthn

import (
	"encoding/json"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

func LoadSessionData(c *gin.Context, username string) (*webauthn.SessionData, error) {
	session := sessions.Default(c)
	// Get session data bytes from the session
	bytes := session.Get(username).([]byte)
	// Unmarshal session data from JSON
	var sessionData webauthn.SessionData
	if err := json.Unmarshal(bytes, &sessionData); err != nil {
		return nil, err
	}
	return &sessionData, nil
}

func StoreSessionData(c *gin.Context, username string, sessionData *webauthn.SessionData) error {
	session := sessions.Default(c)
	// Marshal session data to JSON
	bytes, err := json.Marshal(sessionData)
	if err != nil {
		return fmt.Errorf("failed to marhsall session data: %v", err)	
	}

	// set session data for the user
	session.Set(username, bytes)
	err = session.Save()
	if err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}
	return nil
}
