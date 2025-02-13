package auth

import (
	"net/http"

	"github.com/baala3/passkey-demo/users"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

type WebAuthnController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserStore users.UserRepository
}

func (wc *WebAuthnController) BeginRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		username := ctx.Param("username")

	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)

	if err != nil {
		user, err = wc.UserStore.CreateUser(ctx.Request().Context(), username)
		if err != nil {
			ctx.Logger().Error("error CreateUser() %v", err)
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := wc.WebAuthnAPI.BeginRegistration(user,
		webauthn.WithExclusions(user.CredentialExcludeList()))

	if err != nil{
		ctx.Logger().Error("error webauthnAPI.BeginRegistration() %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := CreateSession(ctx.Request().Context(), username, sessionData); err != nil {
		ctx.Logger().Error("error CreateSession() %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishRegistration() echo.HandlerFunc {
	return func(ctx echo.Context) error {
	username := ctx.Param("username")

	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)

	if err != nil {
		ctx.Logger().Error("error FindUserByName() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	sessionData, err := GetSession(ctx.Request().Context(), username)
	if err != nil {
		ctx.Logger().Error("error GetSession() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	credential, err := wc.WebAuthnAPI.FinishRegistration(user, *sessionData, ctx.Request())
	if err != nil {
		ctx.Logger().Error("error FinishRegistration() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := wc.UserStore.AddWebauthnCredential(ctx.Request().Context(), user.ID, credential); err != nil {
		ctx.Logger().Error("error AddWebauthnCredential() %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
}

func (wc *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
	username := ctx.Param("username")
	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)
	
	if err != nil {
		ctx.Logger().Error("error FindUserByName() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	options, sessionData, err := wc.WebAuthnAPI.BeginLogin(user)
	if err != nil {
		ctx.Logger().Error("error webauthnAPI.BeginLogin() %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := CreateSession(ctx.Request().Context(), username, sessionData); err != nil {
		ctx.Logger().Error("error CreateSession() %v", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, options)
}
}

func (wc *WebAuthnController) FinishLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		username := ctx.Param("username")

	user, err := wc.UserStore.FindUserByName(ctx.Request().Context(), username)
	if err != nil {
		ctx.Logger().Error("error FindUserByName() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	sessionData, err := GetSession(ctx.Request().Context(), username)
	if err != nil {
		ctx.Logger().Error("error GetSession() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	// in an actual implementation we should perform additional
	// checks on the returned 'credential'
	_, err = wc.WebAuthnAPI.FinishLogin(user, *sessionData, ctx.Request())
	if err != nil {
		ctx.Logger().Error("error webauthnAPI.FinishLogin() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
}