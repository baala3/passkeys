package model

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// https://github.com/go-webauthn/webauthn/blob/master/webauthn/credential.go
type WebauthnCredentials struct {
	ID uuid.UUID `json:"id" bun:"id,pk"`
	UserID uuid.UUID `json:"user_id" bun:"user_id"`
	CredentialID []byte `json:"credential_id" bun:"credential_id"`
	PublicKey []byte `json:"public_key" bun:"public_key"`
	AttestationType string `json:"attestation_type" bun:"attestation_type"`
	Transport []protocol.AuthenticatorTransport `json:"transport" bun:"transport,array"`
	Flags webauthn.CredentialFlags `json:"flags" bun:"flags"`
	Authenticator webauthn.Authenticator `json:"authenticator" bun:"authenticator"`
}