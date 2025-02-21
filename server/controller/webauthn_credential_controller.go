package controller

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/baala3/passkeys/concerns"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type WebAuthnCredentialController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserRepository repository.UserRepository
	WebAuthnSession pkg.WebAuthnSession
	UserSession pkg.UserSession
}

func (pc *WebAuthnCredentialController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p pkg.Params
		if err := ctx.Bind(&p); err != nil {
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		if (p == pkg.Params{}) {
			user := concerns.CurrentUser(ctx, pc.UserRepository)
			if user == nil {
				return pkg.SendError(ctx, errors.New("no user found"), http.StatusBadRequest)
			}
			p = pkg.Params{
				Email: user.Email,
			}
		}

		user, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

		if user == nil {
			passwordHash, err := argon2id.CreateHash(random.String(20), argon2id.DefaultParams)
			if err != nil {
				return pkg.SendError(ctx, errors.New("Internal server error"), http.StatusInternalServerError)
			}
			user, err = pc.UserRepository.CreateUser(ctx.Request().Context(), p.Email, passwordHash)
			if err != nil {
				return pkg.SendError(ctx, err, http.StatusInternalServerError)
			}
		}

		authSelect := protocol.AuthenticatorSelection{
			RequireResidentKey: protocol.ResidentKeyRequired(),
			ResidentKey:        protocol.ResidentKeyRequirementRequired,
			UserVerification:   protocol.VerificationRequired,
		}

		// generate PublicKeyCredentialCreationOptions, session data
		options, sessionData, err := pc.WebAuthnAPI.BeginRegistration(user,
			webauthn.WithAuthenticatorSelection(authSelect),
			webauthn.WithExclusions(user.CredentialExcludeList()))

		if err != nil{
			_ = pc.UserRepository.DeleteUser(ctx.Request().Context(), user)
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		err = pc.WebAuthnSession.Create(ctx,"registration", sessionData)
		if err != nil {
			_ = pc.UserRepository.DeleteUser(ctx.Request().Context(), user)
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}
		
		return ctx.JSON(http.StatusOK, options)
	}
}

func (pc *WebAuthnCredentialController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := pc.WebAuthnSession.Get(ctx,"registration")

		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		user, err := pc.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		credential, err := pc.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
		if err != nil {
			_ = pc.UserRepository.DeleteUser(ctx.Request().Context(), user)
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
			_ = pc.UserRepository.DeleteUser(ctx.Request().Context(), user)
			return pkg.SendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
		}

		if err := pc.UserRepository.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
			_ = pc.UserRepository.DeleteUser(ctx.Request().Context(), user)
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		_ = pc.WebAuthnSession.Delete(ctx, sessionId)

		if err := pc.UserSession.Create(ctx, user.ID); err != nil {
			_ = pc.UserRepository.DeleteUser(ctx.Request().Context(), user)
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}
		return pkg.SendOK(ctx)
	}
}

func (pc *WebAuthnCredentialController) GetCredentials() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		user := concerns.CurrentUser(ctx, pc.UserRepository)
		if user == nil {
			return pkg.SendError(ctx, errors.New("user not found"), http.StatusUnauthorized)
		}
		credentials := user.GetWebAuthnCredentials()
		return ctx.JSON(http.StatusOK, credentials)
	}
}

func (pc *WebAuthnCredentialController) DeleteCredential() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		type Request struct {
			CredentialID string `json:"credentialId"` 
		}
		var request Request

		if err := ctx.Bind(&request); err != nil {
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		// Decode the base64 credential ID back to bytes
		credentialID, err := base64.StdEncoding.DecodeString(request.CredentialID)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		user := concerns.CurrentUser(ctx, pc.UserRepository)
		if user == nil {
			return pkg.SendError(ctx, errors.New("user not found"), http.StatusUnauthorized)
		}
		err = pc.UserRepository.DeleteWebauthnCredential(ctx.Request().Context(),user.ID, credentialID)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}
		return pkg.SendOK(ctx)
	}
}
