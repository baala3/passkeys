package model

import (
	"strings"
	"time"

	"github.com/baala3/passkeys/pkg"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id" bun:"id,pk"`
	Email string `json:"email" bun:"email"`
	PasswordHash string `json:"-" bun:"password_hash,notnull"`
	WebauthnCredentials []WebauthnCredentials `json:"webauthn_credentials" bun:"rel:has-many,join:id=user_id"`
	CreatedAt time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at"`
}

// WebAuthnID returns the user's id
func (u *User) WebAuthnID() []byte {
	bytes, _ := u.ID.MarshalBinary()
	return bytes
}

// WebAuthnName returns the user's username
func (u *User) WebAuthnName() string {
	return strings.Split(u.Email, "@")[0]
}

// WebAuthnDisplayName returns the user's display name
func (u *User) WebAuthnDisplayName() string {
	return strings.Split(u.Email, "@")[0]
}

// WebAuthnIcon is not (yet) implemented
func (u *User) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials returns the user's credentials
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	if u.WebauthnCredentials == nil {
		return nil
	}
	credentials := make([]webauthn.Credential, len(u.WebauthnCredentials))
	for i, cred := range u.WebauthnCredentials {
		credentials[i] = webauthn.Credential{
			ID: cred.CredentialID,
			PublicKey: cred.PublicKey,
			AttestationType: cred.AttestationType,
			Transport: cred.Transport,
			Flags: cred.Flags,
		}
	}
	return credentials
}

func (u *User) GetWebAuthnCredentials() []pkg.WebAuthnCredentials {
	credentials := make([]pkg.WebAuthnCredentials, len(u.WebauthnCredentials))
	if err := pkg.LoadAAGUIDs(); err != nil {
		return nil
	}
	for i, cred := range u.WebauthnCredentials {
		credentials[i] = pkg.WebAuthnCredentials{
			PasskeyProvider: pkg.GetPasskeyProviderByAAGUID(cred.Authenticator.AAGUID),
			CredentialId: cred.CredentialID,
			CreatedAt: cred.CreatedAt,
			UpdatedAt: cred.UpdatedAt, // TODO: Add last used at
		}
	}
	return credentials
}

// Returns authenticators already registered to the user
// to prevent multiple registrations of the same authenticator
func (u *User) CredentialExcludeList() []protocol.CredentialDescriptor {
	var credentialExcludeList []protocol.CredentialDescriptor

	for _, cred := range u.WebauthnCredentials {
		descriptor := protocol.CredentialDescriptor{
			Type: protocol.PublicKeyCredentialType,
			CredentialID: cred.CredentialID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}
	return credentialExcludeList
}
