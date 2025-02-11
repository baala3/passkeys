package auth

import "github.com/go-webauthn/webauthn/webauthn"

func NewWebAuthnAPI() (*webauthn.WebAuthn, error) {
	webauthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Passkey Demo",
		RPID: "localhost",
		RPOrigins: []string{"http://localhost:8080"},
	})
	return webauthn, err
}
