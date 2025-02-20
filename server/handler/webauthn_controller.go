package handler

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/baala3/passkeys/repository"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type WebAuthnController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserRepository repository.UserRepository
	SessionRepository repository.SessionRepository
}

type FIDO2Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

func (wc *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil{
			return sendError(ctx, err, http.StatusBadRequest)
		}

		if !validEmail(p.Email){
			return sendError(ctx, errors.New("invalid email"), http.StatusBadRequest)
		}

	_, err := wc.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err == nil {
		return sendError(ctx, errors.New("An account with this email already exists"), http.StatusConflict)
	}

	passwordHash, err := argon2id.CreateHash(random.String(20), argon2id.DefaultParams)
	if err != nil {
		return sendError(ctx, errors.New("Internal server error"), http.StatusInternalServerError)
	}

	user, err := wc.UserRepository.CreateUser(ctx.Request().Context(), p.Email, passwordHash)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	authSelect := protocol.AuthenticatorSelection{
		RequireResidentKey: protocol.ResidentKeyRequired(),
		ResidentKey:        protocol.ResidentKeyRequirementRequired,
		UserVerification:   protocol.VerificationRequired,
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := wc.WebAuthnAPI.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithExclusions(user.CredentialExcludeList()))

	if err != nil{
		_ = wc.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	err = wc.SessionRepository.CreateWebauthnSession(ctx,"registration", sessionData)
	if err != nil {
		_ = wc.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}
	
	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := wc.SessionRepository.GetWebauthnSession(ctx,"registration")

	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	user, err := wc.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	credential, err := wc.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
	if err != nil {
		_ = wc.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		_ = wc.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
	}

	if err := wc.UserRepository.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
		_ = wc.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	_ = wc.SessionRepository.DeleteSession(ctx.Request().Context(), sessionId)

	if err := wc.SessionRepository.Login(ctx, user.ID); err != nil {
		_ = wc.UserRepository.DeleteUser(ctx.Request().Context(), user)
		return sendError(ctx, err, http.StatusInternalServerError)
	}
	return sendOK(ctx)
}
}

func (wc *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return wc.assertionOptions(wc.getCredentialAssertion)
}

func (wc *WebAuthnController) FinishLogin() echo.HandlerFunc {
	return wc.assertionResult(wc.getCredential)
}

func (wc *WebAuthnController) BeginDiscoverableLogin() echo.HandlerFunc {
	return wc.assertionOptions(wc.getDiscoverableCredentialAssertion)
}

func (wc *WebAuthnController) FinishDiscoverableLogin() echo.HandlerFunc {
	return wc.assertionResult(wc.getDiscoverableCredential)
}

func (wc *WebAuthnController) assertionOptions(getCredentialAssertion func(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		options, sessionData, err := getCredentialAssertion(ctx)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		if err := wc.SessionRepository.CreateWebauthnSession(ctx, "login", sessionData); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (wc *WebAuthnController) getCredentialAssertion(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	var p Params
	if err := ctx.Bind(&p); err != nil {
		return nil, nil, err
	}

	user, err := wc.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("User does not exist")
	}

	return wc.WebAuthnAPI.BeginLogin(user, webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (wc *WebAuthnController) assertionResult(getCredential func(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := wc.SessionRepository.GetWebauthnSession(ctx, "login")
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
		wc.SessionRepository.DeleteSession(ctx.Request().Context(), sessionId)

		userID, err := wc.UserRepository.FindUserIDByCredentialID(ctx.Request().Context(), credential.ID)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		if err := wc.SessionRepository.Login(ctx, *userID); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		return sendOK(ctx)
	}
}

func (wc *WebAuthnController) getCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
   user, err := wc.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
   if err != nil {
	return nil, err
   }

   return wc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
}

func (wc *WebAuthnController) getDiscoverableCredentialAssertion(echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	return wc.WebAuthnAPI.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (wc *WebAuthnController) getDiscoverableCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
	return wc.WebAuthnAPI.FinishDiscoverableLogin(
		func(rawId []byte, userID []byte) (user webauthn.User, err error) {
			return wc.UserRepository.FindUserById(ctx.Request().Context(), userID)
		}, *sessionData, ctx.Request())
}
