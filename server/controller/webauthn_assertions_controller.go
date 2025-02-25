package controller

import (
	"errors"
	"net/http"

	"github.com/baala3/passkeys/concerns"
	"github.com/baala3/passkeys/model"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type WebAuthnAssertionsController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserRepository repository.UserRepository
	WebAuthnSession pkg.WebAuthnSession
	UserSession pkg.UserSession
}

func (pc *WebAuthnAssertionsController) BeginLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		options, sessionData, err := pc.getCredentialAssertion(ctx)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if err := pc.WebAuthnSession.Create(ctx, "login", sessionData); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (pc *WebAuthnAssertionsController) FinishLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := pc.WebAuthnSession.Get(ctx, "login")
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		credential, err := pc.getCredential(ctx, sessionData)
		if err != nil {
			return pkg.SendError(ctx, errors.New("There is no password for this account"), http.StatusBadRequest)
		}

		if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
			return pkg.SendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
		}

		if credential.Authenticator.CloneWarning {
			return pkg.SendError(ctx, errors.New("Authenticator is cloned"), http.StatusBadRequest)
		}
		pc.WebAuthnSession.Delete(ctx, sessionId)

		userID, err := pc.UserRepository.FindUserIDByCredentialID(ctx.Request().Context(), credential.ID)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if err := pc.UserSession.Create(ctx, *userID); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		return pkg.SendOK(ctx)
	}
}

func (pc *WebAuthnAssertionsController) BeginDiscoverableLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		options, sessionData, err := pc.getDiscoverableCredentialAssertion()
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if err := pc.WebAuthnSession.Create(ctx, "login", sessionData); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (pc *WebAuthnAssertionsController) FinishDiscoverableLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := pc.WebAuthnSession.Get(ctx, "login")
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		credential, err := pc.getDiscoverableCredential(ctx, sessionData)
		if err != nil {
			return pkg.SendError(ctx, errors.New("There is no password for this account"), http.StatusBadRequest)
		}

		if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
			return pkg.SendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
		}

		if credential.Authenticator.CloneWarning {
			return pkg.SendError(ctx, errors.New("Authenticator is cloned"), http.StatusBadRequest)
		}
		pc.WebAuthnSession.Delete(ctx, sessionId)

		userID, err := pc.UserRepository.FindUserIDByCredentialID(ctx.Request().Context(), credential.ID)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if err := pc.UserSession.Create(ctx, *userID); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		return pkg.SendOK(ctx)
	}
}

func (pc *WebAuthnAssertionsController) getCredentialAssertion(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	user, err, _ := pc.getContextBasedUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	return pc.WebAuthnAPI.BeginLogin(user, webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (pc *WebAuthnAssertionsController) getCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
   user, err := pc.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
   if err != nil {
	return nil, err
   }

   return pc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
}

func (pc *WebAuthnAssertionsController) getDiscoverableCredentialAssertion() (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	return pc.WebAuthnAPI.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (pc *WebAuthnAssertionsController) getDiscoverableCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
	return pc.WebAuthnAPI.FinishDiscoverableLogin(
		func(rawId []byte, userID []byte) (user webauthn.User, err error) {
			return pc.UserRepository.FindUserById(ctx.Request().Context(), userID)
		}, *sessionData, ctx.Request())
}


func (pc *WebAuthnAssertionsController) getContextBasedUser(ctx echo.Context) (user *model.User, err error, status int) {
	context := ctx.QueryParam("context")

	switch context {
	case "signin":
		var p pkg.Params
		if err := ctx.Bind(&p); err != nil {
			return nil, err, http.StatusBadRequest
		}
		user, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

		if err != nil {
			return nil, err, http.StatusBadRequest
		}

		if user == nil {
			return nil, errors.New("User does not exist"), http.StatusBadRequest
		}
		return user, nil, http.StatusOK
	case "delete_account", "email_change", "password_change":
		user = concerns.CurrentUser(ctx, pc.UserRepository)
		if user == nil {
			return nil, errors.New("user not found"), http.StatusUnauthorized
		}
		return user, nil, http.StatusOK
	default:
		return nil, errors.New("invalid context"), http.StatusBadRequest
	}
}
