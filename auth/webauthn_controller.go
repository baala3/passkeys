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

	user, err := wc.UserStore.GetUser(username)

	if err != nil {
		wc.UserStore.PutUser(username)
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := wc.WebAuthnAPI.BeginRegistration(user)

	if err!=nil{
		ctx.Logger().Error("error BeginRegistration() %v", err)
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

	user, err := wc.UserStore.GetUser(username)

	if err != nil {
		ctx.Logger().Error("error GetUser() %v", err)
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

	user.AddCredential(*credential)

	return ctx.JSON(http.StatusOK, nil)
}
}

func (wc *WebAuthnController) BeginLogin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		username := ctx.Param("username")
	 user, err := wc.UserStore.GetUser(username)
	
	if err != nil {
		ctx.Logger().Error("error GetUser() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	options, sessionData, err := wc.WebAuthnAPI.BeginLogin(user)
	if err != nil {
		ctx.Logger().Error("error BeginLogin() %v", err)
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

	user, err := wc.UserStore.GetUser(username)
	if err != nil {
		ctx.Logger().Error("error GetUser() %v", err)
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
		ctx.Logger().Error("error FinishLogin() %v", err)
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
}