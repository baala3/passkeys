//go:build wireinject
// +build wireinject

package main

import (
	"github.com/baala3/passkeys/controller"
	"github.com/baala3/passkeys/db"
	"github.com/baala3/passkeys/pkg"
	"github.com/baala3/passkeys/repository"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

// run `wire` to generate server and all dependencies
func NewServer() (*Server, error) {
	panic(wire.Build(
		wire.Struct(new(Server), "*"),
		echo.New,
		wire.Struct(new(repository.UserRepository), "*"),
		db.GetDB,
		wire.Struct(new(controller.WebAuthnAssertionsController), "*"),
		wire.Struct(new(controller.WebAuthnCredentialController), "*"),
		wire.Struct(new(controller.PasswordController), "*"),
		wire.Struct(new(controller.EmailController), "*"),
		pkg.NewWebAuthnAPI,
		pkg.GetRedisClient,
		wire.Struct(new(pkg.UserSession), "*"),
		wire.Struct(new(pkg.WebAuthnSession), "*"),
	))
}
