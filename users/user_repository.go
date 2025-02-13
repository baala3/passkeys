package users

import (
	"context"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	DB *bun.DB
}

// FindUserByName returns a user by name
func (ur *UserRepository) FindUserByName(ctx context.Context, name string) (*User, error) {
	var user User
	err := ur.DB.NewSelect().
		Model(&user).
		Relation("WebauthnCredentials").
		Column("*").
		Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, name string) (*User, error) {
	user := &User{
		Name: name,
	}

	_, err := ur.DB.NewInsert().
		Model(user).
		Column("name").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) AddWebauthnCredential(ctx context.Context, userID uuid.UUID, credential *webauthn.Credential) error {
	transport := []protocol.AuthenticatorTransport{}
	for _, t := range credential.Transport {
		transport = append(transport, t)
	}

	newWebautnCredential := &WebauthnCredentials{
		UserID: userID,
		CredentialID: credential.ID,
		PublicKey: credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport: credential.Transport,
		Flags: credential.Flags,
		Authenticator: credential.Authenticator,
	}

	_, err := ur.DB.NewInsert().
		Model(newWebautnCredential).
		Column("user_id", "credential_id", "public_key", "attestation_type", "flags", "authenticator").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
