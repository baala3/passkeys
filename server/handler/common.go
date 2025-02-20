package handler

import (
	"net/http"
	"net/mail"

	"github.com/labstack/echo/v4"
)

type Params struct {
	Email string
	Password string
}


type Response struct {
	Status string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

func sendError(ctx echo.Context, err error, code int) error {
	ctx.Logger().Error("Error: %v", err)
	return ctx.JSON(code, Response{
		Status:       "error",
		ErrorMessage: err.Error(),
	})
}

func sendOK(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, Response{
		Status:       "ok",
		ErrorMessage: "",
	})
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
