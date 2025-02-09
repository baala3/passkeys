package main

import (
	"crypto/rand"
	"encoding/binary"

	"github.com/go-webauthn/webauthn/webauthn"
)

type User struct {
	id uint64
	name string
	displayName string
	credentials []webauthn.Credential
}

// NewUser creates a new user with a random id
func NewUser(name, displayName string) *User {
	return &User{
		id: randomUint64(),
		name: name,
		displayName: displayName,
		// user.credentials = []webauthn.Credential{}
	}
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// WebAuthnName returns the user's username
func (u *User) WebAuthnName() string {
	return u.name
}
// WebAuthnDisplayName returns the user's display name
func (u *User) WebAuthnDisplayName() string {
	return u.displayName
}

// WebAuthnID returns the user's id
func (u *User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, u.id)
	return buf
}

// WebAuthnIcon is not (yet) implemented
func (u *User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns the user's credentials
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}
