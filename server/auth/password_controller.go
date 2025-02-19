package auth

import (
	"errors"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/baala3/passkeys/users"
	"github.com/labstack/echo/v4"
)

type PasswordController struct {
	UserRepository users.UserRepository
}

func (pc PasswordController) SignUp() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil {
			return sendError(ctx, err, http.StatusBadRequest)
		}

		email := p.Email
		password := p.Password

		if !validEmail(email) {
			return sendError(ctx, errors.New("invalid email"), http.StatusBadRequest)
		}

		if len(password) < 8 {
			return sendError(ctx, errors.New("password must be at least 8 characters"), http.StatusBadRequest)
		}

		_, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), email)

		if err == nil {
			return sendError(ctx, errors.New("An account with this email already exists"), http.StatusBadRequest)
		}

		passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}

		_, err = pc.UserRepository.CreateUser(ctx.Request().Context(), email, passwordHash)
		if err != nil {
			return sendError(ctx, err, http.StatusInternalServerError)
		}
		return sendOK(ctx)
		
	}
}

func (pc PasswordController) Login() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p Params
		if err := ctx.Bind(&p); err != nil {
			return sendError(ctx, err, http.StatusBadRequest)
		}

		email := p.Email

		if !validEmail(email) {
			return sendError(ctx, errors.New("invalid email"), http.StatusBadRequest)
		}

		user, err := pc.UserRepository.FindUserByEmail(ctx.Request().Context(), email)
		if err != nil {
			return sendError(ctx, errors.New("An account with this email does not exist"), http.StatusBadRequest)
		}

		match, err := argon2id.ComparePasswordAndHash(p.Password, user.PasswordHash)
		if err != nil {
			return sendError(ctx, err, http.StatusBadRequest)
		}

		if !match {
			return sendError(ctx, errors.New("Invalid password"), http.StatusBadRequest)
		}
		return sendOK(ctx)
	}
}
