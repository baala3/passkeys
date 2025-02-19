package auth

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/baala3/passkeys/users"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type WebAuthnController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserStore users.UserRepository
}

type FIDO2Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

func (wc *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p AuthParams
		if err := ctx.Bind(&p); err != nil{
			return sendError(ctx, err, http.StatusBadRequest)
		}

		if !validEmail(p.Email){
			return sendError(ctx, errors.New("invalid email"), http.StatusBadRequest)
		}

	_, err := wc.UserStore.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err == nil {
		return sendError(ctx, errors.New("An account with this email already exists"), http.StatusConflict)
	}

	passwordHash, err := argon2id.CreateHash(random.String(20), argon2id.DefaultParams)
	if err != nil {
		return sendError(ctx, errors.New("Internal server error"), http.StatusInternalServerError)
	}

	user, err := wc.UserStore.CreateUser(ctx.Request().Context(), p.Email, passwordHash)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	authSelect := protocol.AuthenticatorSelection{
		RequireResidentKey: protocol.ResidentKeyRequired(),
		ResidentKey:        protocol.ResidentKeyRequirementRequired,
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := wc.WebAuthnAPI.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithExclusions(user.CredentialExcludeList()))

	if err != nil{
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	err = CreateSession(ctx,"registration", sessionData)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}
	
	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := GetSession(ctx,"registration")

	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	user, err := wc.UserStore.FindUserById(ctx.Request().Context(), sessionData.UserID)
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	credential, err := wc.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
	if err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		return sendError(ctx, errors.New("User not present or not verified"), http.StatusBadRequest)
	}

	if err := wc.UserStore.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
		return sendError(ctx, err, http.StatusInternalServerError)
	}

	DeleteSession(ctx.Request().Context(), sessionId)
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

		if err := CreateSession(ctx, "login", sessionData); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (wc *WebAuthnController) getCredentialAssertion(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	var p AuthParams
	if err := ctx.Bind(&p); err != nil {
		return nil, nil, err
	}

	user, err := wc.UserStore.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("User does not exist")
	}

	return wc.WebAuthnAPI.BeginLogin(user)
}

func (wc *WebAuthnController) assertionResult(getCredential func(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := GetSession(ctx, "login")
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
		DeleteSession(ctx.Request().Context(), sessionId)
		return sendOK(ctx)
	}
}

func (wc *WebAuthnController) getCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
   user, err := wc.UserStore.FindUserById(ctx.Request().Context(), sessionData.UserID)
   if err != nil {
	return nil, err
   }

   return wc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
}

func (wc *WebAuthnController) getDiscoverableCredentialAssertion(echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	return wc.WebAuthnAPI.BeginDiscoverableLogin()
}

func (wc *WebAuthnController) getDiscoverableCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
	return wc.WebAuthnAPI.FinishDiscoverableLogin(
		func(rawId []byte, userID []byte) (user webauthn.User, err error) {
			return wc.UserStore.FindUserById(ctx.Request().Context(), userID)
		}, *sessionData, ctx.Request())
}
