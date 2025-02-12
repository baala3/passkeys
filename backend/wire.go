//go:build wireinject
// +build wireinject

package main

import (
	auth "github.com/baala3/passkey-demo/auth"
	"github.com/baala3/passkey-demo/users"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

func NewServer() (*Server, error) {
	panic(wire.Build(
		wire.Struct(new(Server), "*"),
		echo.New,
		wire.Struct(new(auth.WebAuthnController), "*"),
		users.NewUserRepository,
		auth.NewWebAuthnAPI,
	))
}
