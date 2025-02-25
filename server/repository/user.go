package repository

import (
	"context"

	"github.com/baala3/passkeys/model"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	DB *bun.DB
}

// FindUserByName returns a user by name
func (ur *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := ur.DB.NewSelect().
		Model(&user).
		Relation("WebauthnCredentials").
		Column("*").
		Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserById returns a user by id
func (ur *UserRepository) FindUserById(ctx context.Context, rawUserID []byte) (*model.User, error) {
	userID, err := uuid.FromBytes(rawUserID)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = ur.DB.NewSelect().
		Model(&user).
		Relation("WebauthnCredentials").
		Column("*").
		Where("id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (ur *UserRepository) CreateUser(ctx context.Context, email string, passwordHash string) (*model.User, error) {
	user := &model.User{
		ID: uuid.New(),
		Email: email,
		PasswordHash: passwordHash,
	}

	_, err := ur.DB.NewInsert().
		Model(user).
		Column("id", "email", "password_hash").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserRepository) AddWebauthnCredential(ctx context.Context, userID uuid.UUID, credential *webauthn.Credential) error {
	newWebauthnCredential := &model.WebauthnCredentials{
		ID: uuid.New(),
		UserID: userID,
		CredentialID: credential.ID,
		PublicKey: credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport: credential.Transport,
		Flags: credential.Flags,
		Authenticator: credential.Authenticator,
	}

	_, err := ur.DB.NewInsert().
		Model(newWebauthnCredential).
		Column("id", "user_id", "credential_id", "public_key", "attestation_type", "transport","flags", "authenticator").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) DeleteWebauthnCredential(ctx context.Context, userID uuid.UUID, credentialID []byte) error {
	_, err := ur.DB.NewDelete().
		Model(&model.WebauthnCredentials{}).
		Where("user_id = ?", userID).
		Where("credential_id = ?", credentialID).
		Exec(ctx)
	return err
}

func (ur *UserRepository) DeleteUser(ctx context.Context, user *model.User) error {
	err := ur.DB.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Delete webauthn credentials first (foreign key dependency)
		_, err := tx.NewDelete().
			Model(&model.WebauthnCredentials{}).
			Where("user_id = ?", user.ID).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Delete the user
		_, err = tx.NewDelete().
			Model(user).
			WherePK().
			Exec(ctx)
		return err
	})

	return err
}

func (ur *UserRepository) FindUserIDByCredentialID(ctx context.Context, credentialID []byte) (*uuid.UUID, error) {
	var credential model.WebauthnCredentials
	err := ur.DB.NewSelect().
		Model(&credential).
		Column("user_id").
		Where("credential_id = ?", credentialID).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &credential.UserID, nil
}
