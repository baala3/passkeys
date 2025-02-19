//go:build wireinject
// +build wireinject

package main

import (
	auth "github.com/baala3/passkeys/auth"
	"github.com/baala3/passkeys/db"
	"github.com/baala3/passkeys/users"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

func NewServer() (*Server, error) {
	panic(wire.Build(
		wire.Struct(new(Server), "*"),
		echo.New,
		wire.Struct(new(users.UserRepository), "*"),
		db.GetDB,
		wire.Struct(new(auth.WebAuthnController), "*"),
		auth.NewWebAuthnAPI,
	))
}
