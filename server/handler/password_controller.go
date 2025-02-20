package handler

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

func (handler PasswordController) SignUp() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil {
			return sendError(ctx, err, http.StatusBadRequest)
		}

		email := p.Email
		password := p.Password

		if !validEmail(email) {
			return sendError(ctx, errors.New("Invalid email"), http.StatusBadRequest)
		}

		if len(password) < 8 {
			return sendError(ctx, errors.New("Password must be at least 8 characters"), http.StatusBadRequest)
		}

		_, err := handler.UserRepository.FindUserByEmail(ctx.Request().Context(), email)

		if err == nil {
			return sendError(ctx, errors.New("An account with that email already exists."), http.StatusConflict)
		}

		passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		user, err := handler.UserRepository.CreateUser(ctx.Request().Context(), email, passwordHash)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		if err = handler.UserSession.Create(ctx, user.ID); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}
		return sendOK(ctx)
		
	}
}

func (handler PasswordController) Login() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil {
			return sendError(ctx, err, http.StatusBadRequest)
		}

		email := p.Email

		if !validEmail(email) {
			return sendError(ctx, errors.New("Invalid email"), http.StatusBadRequest)
		}

		user, err := handler.UserRepository.FindUserByEmail(ctx.Request().Context(), email)
		if err != nil {
			return sendError(ctx, errors.New("An account with that email does not exist."), http.StatusNotFound)
		}

		match, err := argon2id.ComparePasswordAndHash(p.Password, user.PasswordHash)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		if !match {
			return sendError(ctx, errors.New("Invalid password."), http.StatusUnauthorized)
		}

		if err = handler.UserSession.Create(ctx, user.ID); err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}
		return sendOK(ctx)
	}
}

func (handler PasswordController) Logout() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("auth")
		if err != nil {
			return sendError(ctx, errors.New("Not logged in."), http.StatusUnauthorized)
		}

		sessionID := cookie.Value
		handler.UserSession.Delete(ctx, sessionID)

		return sendOK(ctx)
	}
}
