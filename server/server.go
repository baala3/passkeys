package main

import (
	"embed"

	"github.com/baala3/passkeys/controller"
	"github.com/baala3/passkeys/middleware"
	"github.com/labstack/echo/v4"
)

type Server struct {
	router             *echo.Echo
	webauthnAssertionsController controller.WebAuthnAssertionsController
	webauthnCredentialController controller.WebAuthnCredentialController
	passwordController controller.PasswordController
}

func (s *Server) Start() {
	s.registerEndpoints()
	s.router.Logger.Fatal(s.router.Start(":9044"))
}

var (
	//go:embed dist/**
	dist embed.FS
	//go:embed dist/index.html
	indexHTML embed.FS

	distDirFS     = echo.MustSubFS(dist, "dist")
	distIndexHTML = echo.MustSubFS(indexHTML, "dist")
)

func (s *Server) registerEndpoints() {
	s.router.Use(middleware.InjectRedis)

	s.router.StaticFS("/", distDirFS)

	s.router.FileFS("/", "index.html", distIndexHTML)
	s.router.FileFS("/sign-up", "index.html", distIndexHTML)
	s.router.FileFS("/home", "index.html", distIndexHTML, middleware.Auth)
	s.router.FileFS("/passkeys", "index.html", distIndexHTML, middleware.Auth)

	s.router.POST("/register/begin", s.webauthnCredentialController.BeginRegistration())
	s.router.POST("/register/finish", s.webauthnCredentialController.FinishRegistration())
	s.router.POST("/login/begin", s.webauthnAssertionsController.BeginLogin())
	s.router.POST("/login/finish", s.webauthnAssertionsController.FinishLogin())

	s.router.POST("/discoverable_login/begin", s.webauthnAssertionsController.BeginDiscoverableLogin())
	s.router.POST("/discoverable_login/finish", s.webauthnAssertionsController.FinishDiscoverableLogin())

	s.router.POST("/register/password", s.passwordController.SignUp())
	s.router.POST("/login/password", s.passwordController.Login())
	s.router.POST("/logout", s.passwordController.Logout())
}
