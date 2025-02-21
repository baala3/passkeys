package controller

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/labstack/echo/v4"
)

type PasswordController struct {
	UserRepository repository.UserRepository
	UserSession pkg.UserSession
}

func (pc PasswordController) SignUp() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p pkg.Params
		if err := ctx.Bind(&p); err != nil {
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		email := p.Email
		password := p.Password

		if !pkg.IsValidEmail(email) {
			return pkg.SendError(ctx, errors.New("Invalid email"), http.StatusBadRequest)
		}

		if len(password) < 8 {
			return pkg.SendError(ctx, errors.New("Password must be at least 8 characters"), http.StatusBadRequest)
		}

		_, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), email)

		if err == nil {
			return pkg.SendError(ctx, errors.New("An account with that email already exists."), http.StatusConflict)
		}

		passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		user, err := pc.UserRepository.CreateUser(ctx.Request().Context(), email, passwordHash)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if err = pc.UserSession.Create(ctx, user.ID); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}
		return pkg.SendOK(ctx)
		
	}
}

func (pc PasswordController) Login() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p pkg.Params
		if err := ctx.Bind(&p); err != nil {
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		email := p.Email

		if !pkg.IsValidEmail(email) {
			return pkg.SendError(ctx, errors.New("Invalid email"), http.StatusBadRequest)
		}

		user, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), email)
		if err != nil {
			return pkg.SendError(ctx, errors.New("An account with that email does not exist."), http.StatusNotFound)
		}

		match, err := argon2id.ComparePasswordAndHash(p.Password, user.PasswordHash)
		if err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		if !match {
			return pkg.SendError(ctx, errors.New("Invalid password."), http.StatusUnauthorized)
		}

		if err = pc.UserSession.Create(ctx, user.ID); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}
		return pkg.SendOK(ctx)
	}
}

func (pc PasswordController) Logout() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("auth")
		if err != nil {
			return pkg.SendError(ctx, errors.New("Not logged in."), http.StatusUnauthorized)
		}

		sessionID := cookie.Value
		pc.UserSession.Delete(ctx, sessionID)

		return pkg.SendOK(ctx)
	}
}
