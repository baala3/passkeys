package pkg

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WebAuthnCredentials struct {
	AuthenticatorMetadata AuthenticatorMetadata `json:"authenticator_metadata"`
	CredentialId []byte `json:"credential_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Params struct {
	Email string
	Password string
}

type Response struct {
	Status string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

func SendError(ctx echo.Context, err error, code int) error {
	ctx.Logger().Error("Error: %v", err)
	return ctx.JSON(code, Response{
		Status:       "error",
		ErrorMessage: err.Error(),
	})
}

func SendOK(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Status:       "ok",
		ErrorMessage: "",
	})
}

type AuthenticatorMetadata struct {
    Name      string `json:"name"`
    IconDark  string `json:"icon_dark"`
    IconLight string `json:"icon_light"`
}

var AAGUIDs map[string]AuthenticatorMetadata

func initAAGUIDs() error {
	data, err := os.ReadFile("passkey-authenticator-aaguids/aaguid.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &AAGUIDs)
	if err != nil {
		return err
	}

	return nil
}
func GetPasskeyProviderData(aaguid []byte) AuthenticatorMetadata {
	if AAGUIDs == nil {
		err := initAAGUIDs()
		if err != nil {
			return AuthenticatorMetadata{}
		}
	}

	uuid, err := uuid.FromBytes(aaguid)
	if err != nil {
		return AuthenticatorMetadata{}
	}

	return AAGUIDs[uuid.String()]
}
