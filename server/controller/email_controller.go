package controller

import (
	"errors"
	"net/http"

	"github.com/baala3/passkeys/concerns"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/labstack/echo/v4"
)

type EmailController struct {
	UserRepository repository.UserRepository
	UserSession pkg.UserSession
}

type EmailParams struct {
	CurrentEmail string 
	NewEmail string
}

func (ec *EmailController) ChangeEmail() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var p EmailParams
		if err := ctx.Bind(&p); err != nil {
			return pkg.SendError(ctx, err, http.StatusBadRequest)
		}

		user := concerns.CurrentUser(ctx, ec.UserRepository)
		if user == nil {
			return pkg.SendError(ctx, errors.New("Cannot get current user"), http.StatusInternalServerError)
		}

		if !pkg.IsValidEmail(p.NewEmail) {
			return pkg.SendError(ctx, errors.New("New email is invalid"), http.StatusBadRequest)
		}

		if user.Email != p.CurrentEmail {
			return pkg.SendError(ctx, errors.New("Current email is incorrect"), http.StatusUnauthorized)
		}

		if user.Email == p.NewEmail {
			return pkg.SendError(ctx, errors.New("New email cannot be the same as the current email"), http.StatusBadRequest)
		}

		existingUser, err := ec.UserRepository.FindUserByEmail(ctx.Request().Context(), p.NewEmail)

		if err == nil || existingUser != nil {
			return pkg.SendError(ctx, errors.New("New email is already in use"), http.StatusBadRequest)
		}

		user.Email = p.NewEmail
		if err := ec.UserRepository.UpdateUser(ctx.Request().Context(), user); err != nil {
			return pkg.SendError(ctx, err, http.StatusInternalServerError)
		}

		return pkg.SendOK(ctx)
	}
}

