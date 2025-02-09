package main

import (
	"crypto/rand"
	"encoding/binary"
)

type User struct {
	id uint64
	name string
	displayName string
	credentials []string
}

// NewUser creates a new user with a random id
func NewUser(name, displayName string, credentials []string) *User {
	return &User{
		id: randomUint64(),
		name: name,
		displayName: displayName,
		credentials: credentials,
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
