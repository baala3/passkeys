package users

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

	// rand.Read generates random bytes and stores them in 'buf'
	_, _ = rand.Read(buf)

	// Convert the byte slice into a uint64 number using LittleEndian format
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
	// Create a byte slice of size 10 (MaxVarintLen64), enough to hold the largest possible Varint (uint64)
	buf := make([]byte, binary.MaxVarintLen64)

	// Encode u.id as a Varint and store it in 'buf'
    // 'PutUvarint' writes the encoded value into the byte slice, 
    // but the byte slice will have a fixed size (10 bytes).
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
