package controller

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

func (pc *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p pkg.Params
		if err := ctx.Bind(&p); err != nil{
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		if !pkg.IsValidEmail(p.Email){
			return pkg.SendError(ctx, errors.New("Invalid email"), http.StatusBadRequest)
		}

	_, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err == nil {
		return pkg.SendError(ctx, errors.New("An account with that email already exists."), http.StatusConflict)
	}

	passwordHash, err := argon2id.CreateHash(random.String(20), argon2id.DefaultParams)
	if err != nil {
		return pkg.SendError(ctx, errors.New("Internal server error"), http.StatusInternalServerError)
	}

	user, err := pc.UserRepository.CreateUser(ctx.Request().Context(), p.Email, passwordHash)
	if err != nil {
		return pkg.SendError(ctx, err, http.StatusInternalServerError)
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

func (pc *WebAuthnController) FinishRegistration() echo.HandlerFunc {
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

func (pc *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return pc.assertionOptions(pc.getCredentialAssertion)
}

func (pc *WebAuthnController) FinishLogin() echo.HandlerFunc {
	return pc.assertionResult(pc.getCredential)
}

func (pc *WebAuthnController) BeginDiscoverableLogin() echo.HandlerFunc {
	return pc.assertionOptions(pc.getDiscoverableCredentialAssertion)
}

func (pc *WebAuthnController) FinishDiscoverableLogin() echo.HandlerFunc {
	return pc.assertionResult(pc.getDiscoverableCredential)
}

func (pc *WebAuthnController) assertionOptions(getCredentialAssertion func(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		options, sessionData, err := getCredentialAssertion(ctx)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if err := pc.WebAuthnSession.Create(ctx, "login", sessionData); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (pc *WebAuthnController) getCredentialAssertion(ctx echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	var p pkg.Params
	if err := ctx.Bind(&p); err != nil {
		return nil, nil, err
	}

	user, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err != nil {
		return nil, nil, err
	}

	if user == nil {
		return nil, nil, errors.New("User does not exist")
	}

	return pc.WebAuthnAPI.BeginLogin(user, webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (pc *WebAuthnController) assertionResult(getCredential func(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error)) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := pc.WebAuthnSession.Get(ctx, "login")
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		credential, err := getCredential(ctx, sessionData)
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

func (pc *WebAuthnController) getCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
   user, err := pc.UserRepository.FindUserById(ctx.Request().Context(), sessionData.UserID)
   if err != nil {
	return nil, err
   }

   return pc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
}

func (pc *WebAuthnController) getDiscoverableCredentialAssertion(echo.Context) (*protocol.CredentialAssertion, *webauthn.SessionData, error) {
	return pc.WebAuthnAPI.BeginDiscoverableLogin(webauthn.WithUserVerification(protocol.VerificationRequired))
}

func (pc *WebAuthnController) getDiscoverableCredential(ctx echo.Context, sessionData *webauthn.SessionData) (*webauthn.Credential, error) {
	return pc.WebAuthnAPI.FinishDiscoverableLogin(
		func(rawId []byte, userID []byte) (user webauthn.User, err error) {
			return pc.UserRepository.FindUserById(ctx.Request().Context(), userID)
		}, *sessionData, ctx.Request())
}
