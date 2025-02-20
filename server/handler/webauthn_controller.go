package handler

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type WebAuthnController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserRepository repository.UserRepository
	WebAuthnSession pkg.WebAuthnSession
	UserSession pkg.UserSession
}

type FIDO2Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

func (handler *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil{
			return sendError(ctx, err, http.StatusBadRequest)
		}

		if !validEmail(p.Email){
			return sendError(ctx, errors.New("Invalid email"), http.StatusBadRequest)
		}

	_, err := handler.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err == nil {
		return sendError(ctx, errors.New("An account with that email already exists."), http.StatusConflict)
	}

	passwordHash, err := argon2id.CreateHash(random.String(20), argon2id.DefaultParams)
	if err != nil {
		return sendError(ctx, errors.New("Internal server error"), http.StatusInternalServerError)
	}

	user, err := handler.UserRepository.CreateUser(ctx.Request().Context(), p.Email, passwordHash)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	authSelect := protocol.AuthenticatorSelection{
		RequireResidentKey: protocol.ResidentKeyRequired(),
		ResidentKey:        protocol.ResidentKeyRequirementRequired,
		UserVerification:   protocol.VerificationRequired,
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := handler.WebAuthnAPI.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithExclusions(user.CredentialExcludeList()))

	if err != nil{
		_ = handler.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	err = handler.WebAuthnSession.Create(ctx,"registration", sessionData)
	if err != nil {
		_ = handler.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}
	
	return ctx.JSON(http.StatusOK, options)
}
}

func (handler *WebAuthnController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := handler.WebAuthnSession.Get(ctx,"registration")

	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	user, err := handler.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	credential, err := handler.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
	if err != nil {
		_ = handler.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		_ = handler.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
	}

	if err := handler.UserRepository.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
		_ = handler.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	_ = handler.WebAuthnSession.Delete(ctx, sessionId)

	if err := handler.UserSession.Create(ctx, user.ID); err != nil {
		_ = handler.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}
	return sendOK(ctx)
}
}

func (handler *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return handler.assertionOptions(handler.getCredentialAssertion)
}

func (handler *WebAuthnController) FinishLogin() echo.HandlerFunc {
	return handler.assertionResult(handler.getCredential)
}

func (handler *WebAuthnController) BeginDiscoverableLogin() echo.HandlerFunc {
	return handler.assertionOptions(handler.getDiscoverableCredentialAssertion)
}

func (handler *WebAuthnController) FinishDiscoverableLogin() echo.HandlerFunc {
	return handler.assertionResult(handler.getDiscoverableCredential)
}

func (handler *WebAuthnController) assertionOptions(getCredentialAssertion func(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		options, sessionData, err := getCredentialAssertion(ctx)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		if err := handler.WebAuthnSession.Create(ctx, "login", sessionData); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (handler *WebAuthnController) getCredentialAssertion(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	var p Params
	if err := ctx.Bind(&p); err != nil {
		return nil, nil, err
	}

	user, err := handler.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("User does not exist")
	}

	return handler.WebAuthnAPI.BeginLogin(user, webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (handler *WebAuthnController) assertionResult(getCredential func(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := handler.WebAuthnSession.Get(ctx, "login")
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		credential, err := getCredential(ctx, sessionData)
		if err != nil {
			return sendError(ctx, errors.New("There is no password for this account"), http.StatusBadRequest)
		}

		if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
			return sendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
		}

		if credential.Authenticator.CloneWarning {
			return sendError(ctx, errors.New("Authenticator is cloned"), http.StatusBadRequest)
		}
		handler.WebAuthnSession.Delete(ctx, sessionId)

		userID, err := handler.UserRepository.FindUserIDByCredentialID(ctx.Request().Context(), credential.ID)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		if err := handler.UserSession.Create(ctx, *userID); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		return sendOK(ctx)
	}
}

func (handler *WebAuthnController) getCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
   user, err := handler.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
   if err != nil {
	return nil, err
   }

   return handler.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
}

func (handler *WebAuthnController) getDiscoverableCredentialAssertion(echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	return handler.WebAuthnAPI.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (handler *WebAuthnController) getDiscoverableCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
	return handler.WebAuthnAPI.FinishDiscoverableLogin(
		func(rawId []byte, userID []byte) (user webauthn.User, err error) {
			return handler.UserRepository.FindUserById(ctx.Request().Context(), userID)
		}, *sessionData, ctx.Request())
}
