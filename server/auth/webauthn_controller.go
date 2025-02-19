package auth

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/baala3/passkey-demo/users"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type WebAuthnController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserStore users.UserRepository
}

type FIDO2Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type Params struct {
	Email string
}

func (wc *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil{
			return sendError(ctx, err)
		}

		if !validEmail(p.Email){
			return sendError(ctx, errors.New("invalid email"))
		}


	user, err := wc.UserStore.FindUserByEmail(ctx.Request().Context(), p.Email)

	if err != nil {
		user, err = wc.UserStore.CreateUser(ctx.Request().Context(), p.Email)
		if err != nil {
			return sendError(ctx, err)
		}
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
		return sendError(ctx, err)
	}

	err = CreateSession(ctx,"registration", sessionData)
	if err != nil {
		return sendError(ctx, err)
	}
	
	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := GetSession(ctx,"registration")

	if err != nil {
		return sendError(ctx, err)
	}

	user, err := wc.UserStore.FindUserById(ctx.Request().Context(), sessionData.UserID)
	if err != nil {
		return sendError(ctx, err)
	}

	credential, err := wc.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
	if err != nil {
		return sendError(ctx, err)
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		return sendError(ctx, errors.New("User not present or not verified"))
	}

	if err := wc.UserStore.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
		return sendError(ctx, err)
	}

	DeleteSession(ctx.Request().Context(), sessionId)
	return sendOK(ctx)
}
}

func (wc *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
	var p Params
	if err := ctx.Bind(&p); err != nil{
		return sendError(ctx, err)
	}

	user, err := wc.UserStore.FindUserByEmail(ctx.Request().Context(), p.Email)
	
	if err != nil {
		return sendError(ctx, err)
	}

	if user == nil {
		return sendError(ctx, errors.New("User does not exist"))
	}

	options, sessionData, err := wc.WebAuthnAPI.BeginLogin(user)
	if err != nil {
		return sendError(ctx, err)
	}

	 err = CreateSession(ctx,"login", sessionData)
	if err != nil {
		return sendError(ctx, err)
	}

	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
	sessionId, sessionData, err := GetSession(ctx,"login")
	
	if err != nil {
		return sendError(ctx, err)
	}

	user, err := wc.UserStore.FindUserById(ctx.Request().Context(), sessionData.UserID)
	if err != nil {
		return sendError(ctx, err)
	}
	credential, err := wc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
	if err != nil {
		return sendError(ctx, err)
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		return sendError(ctx, errors.New("User not present or not verified"))
	}
	if credential.Authenticator.CloneWarning {
		return sendError(ctx, errors.New("Authenticator is cloned"))
	}

	DeleteSession(ctx.Request().Context(), sessionId)
	return sendOK(ctx)
}
}

func (wc *WebAuthnController) BeginDiscoverableLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		options, sessionData, err := wc.WebAuthnAPI.BeginDiscoverableLogin()
		if err != nil {
			return sendError(ctx, err)
		}

		if err := CreateSession(ctx, "login", sessionData); err != nil {
			return sendError(ctx, err)
		}

		return ctx.JSON(http.StatusOK, options)
	}
}

func (wc *WebAuthnController) FinishDiscoverableLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sessionId, sessionData, err := GetSession(ctx, "login")
		if err != nil {
			return sendError(ctx, err)
		}
		credential, err := wc.WebAuthnAPI.FinishDiscoverableLogin(
			func(rawId []byte, userID []byte) (user webauthn.User, err error) {
				return wc.UserStore.FindUserById(ctx.Request().Context(), userID)
				}, *sessionData, ctx.Request())

			if err != nil {
				return sendError(ctx, err)
			}

			if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
				return sendError(ctx, errors.New("User not present or not verified"))
			} 

			if credential.Authenticator.CloneWarning {
				return sendError(ctx, errors.New("Authenticator is cloned"))
			}

			DeleteSession(ctx.Request().Context(), sessionId)
			return sendOK(ctx)
	}
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func sendError(ctx echo.Context, err error) error {
	ctx.Logger().Error("Error: %v", err)
	return ctx.JSON(http.StatusBadRequest, FIDO2Response{
		Status:       "error",
		ErrorMessage: err.Error(),
	})
}

func sendOK(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, FIDO2Response{
		Status:       "ok",
		ErrorMessage: "",
	})
}
