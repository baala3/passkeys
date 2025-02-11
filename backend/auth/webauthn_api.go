package auth

import (
	"os"

	"github.com/go-webauthn/webauthn/webauthn"
)

func NewWebAuthnAPI() (*webauthn.WebAuthn, error) {
	webauthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: os.Getenv("RP_DISPLAY_NAME"),
		RPID: os.Getenv("RP_ID"),
		RPOrigins: []string{os.Getenv("RP_ORIGIN")},
	})
	return webauthn, err
}
