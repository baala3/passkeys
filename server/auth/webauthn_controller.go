package auth

import (
	"net/http"

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

func (wc *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		username := ctx.Param("username")

	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)

	if err != nil {
		user, err = wc.UserStore.CreateUser(ctx.Request().Context(), username)
		if err != nil {
			ctx.Logger().Error("error CreateUser() %v", err)
			return ctx.JSON(http.StatusInternalServerError, FIDO2Response{
				Status:       "error",
				ErrorMessage: err.Error(),
			})
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
		ctx.Logger().Error("error webauthnAPI.BeginRegistration() %v", err)
		return ctx.JSON(http.StatusInternalServerError, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	sessionId, err := CreateSession(ctx.Request().Context(), sessionData)
	if err != nil {
		ctx.Logger().Error("error CreateSession() %v", err)
		return ctx.JSON(http.StatusInternalServerError, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	ctx.SetCookie(&http.Cookie{
		Name: "registration",
		Value: sessionId,
		Path: "/",
	})

	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
	username := ctx.Param("username")

	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)

	if err != nil {
		ctx.Logger().Error("error FindUserByName() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	cookie, err := ctx.Cookie("registration")
	if err != nil {
		ctx.Logger().Error("error GetCookie() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	sessionId := cookie.Value
	sessionData, err := GetSession(ctx.Request().Context(), sessionId)
	if err != nil {
		ctx.Logger().Error("error GetSession() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	credential, err := wc.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
	if err != nil {
		ctx.Logger().Error("error FinishRegistration() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: "User not present or not verified",
		})
	}

	if err := wc.UserStore.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
		ctx.Logger().Error("error AddWebauthnCredential() %v", err)
		return ctx.JSON(http.StatusInternalServerError, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	DeleteSession(ctx.Request().Context(), sessionId)
	return ctx.JSON(http.StatusOK, FIDO2Response{
		Status:       "ok",
		ErrorMessage: "",
	})
}
}

func (wc *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
	username := ctx.Param("username")
	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)
	
	if err != nil {
		ctx.Logger().Error("error FindUserByName() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	if user == nil {
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: "User does not exist",
		})
	}

	options, sessionData, err := wc.WebAuthnAPI.BeginLogin(user)
	if err != nil {
		ctx.Logger().Error("error webauthnAPI.BeginLogin() %v", err)
		return ctx.JSON(http.StatusInternalServerError, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	sessionId, err := CreateSession(ctx.Request().Context(), sessionData)
	if err != nil {
		ctx.Logger().Error("error CreateSession() %v", err)
		return ctx.JSON(http.StatusInternalServerError, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	ctx.SetCookie(&http.Cookie{
		Name: "login",
		Value: sessionId,
		Path: "/",
	})

	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		username := ctx.Param("username")

	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)
	if err != nil {
		ctx.Logger().Error("error FindUserByName() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	cookie, err := ctx.Cookie("login")
	sessionId := cookie.Value
	if err != nil {
		ctx.Logger().Error("error GetCookie() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}
	sessionData, err := GetSession(ctx.Request().Context(), sessionId)
	if err != nil {
		ctx.Logger().Error("error GetSession() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	credential, err := wc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
	if err != nil {
		ctx.Logger().Error("error webauthnAPI.FinishLogin() %v", err)
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: err.Error(),
		})
	}

	if !credential.Flags.UserPresent || !credential.Flags.UserVerified {
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: "User not present or not verified",
		})
	}
	if credential.Authenticator.CloneWarning {
		return ctx.JSON(http.StatusBadRequest, FIDO2Response{
			Status:       "error",
			ErrorMessage: "Authenticator is cloned",
		})
	}

	DeleteSession(ctx.Request().Context(), sessionId)
	return ctx.JSON(http.StatusOK, FIDO2Response{
		Status:       "ok",
		ErrorMessage: "",
	})
}
}