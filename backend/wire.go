//go:build wireinject
// +build wireinject

package main

import (
	auth "github.com/baala3/passkey-demo/auth"
	"github.com/baala3/passkey-demo/users"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func provideGinOptions() []gin.OptionFunc {
	return []gin.OptionFunc{}
}

var engineSet = wire.NewSet(
	gin.Default,
	provideGinOptions,
)

func NewServer() (*Server, error) {
	panic(wire.Build(
		wire.Struct(new(Server), "*"),
		engineSet,
		wire.Struct(new(auth.WebAuthnController), "*"),
		users.NewUserRepository,
		auth.NewWebAuthnAPI,
	))
}
