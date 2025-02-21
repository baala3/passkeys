package pkg

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type WebAuthnCredentials struct {
	AAGUID []byte `json:"aaguid" bun:"aaguid"`
	CredentialId []byte `json:"credential_id" bun:"credential_id"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`
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
