package pkg

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
)

type PasskeyProvider struct {
    Name      string `json:"name"`
    IconDark  string `json:"icon_dark"`
    IconLight string `json:"icon_light"`
}

type WebAuthnCredentials struct {
	PasskeyProvider PasskeyProvider `json:"passkey_provider"`
	CredentialId []byte `json:"credential_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	AAGUIDs   map[string]PasskeyProvider
	loadError error
)

// LoadAAGUIDs should be called explicitly in a specific file
func LoadAAGUIDs() error {
	AAGUIDs = make(map[string]PasskeyProvider)

	data, err := os.ReadFile("passkey-authenticator-aaguids/aaguid.json")
	if err != nil {
		loadError = err
		return err
	}

	if err := json.Unmarshal(data, &AAGUIDs); err != nil {
		loadError = err
		return err
	}

	return nil
}

func GetPasskeyProviderByAAGUID(aaguid []byte) PasskeyProvider {
	if loadError != nil {
		return PasskeyProvider{}
	}

	uuid, err := uuid.FromBytes(aaguid)
	if err != nil {
		return PasskeyProvider{}
	}

	return AAGUIDs[uuid.String()]
}
