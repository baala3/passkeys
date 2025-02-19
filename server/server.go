package main

import (
	"embed"

	"github.com/baala3/passkeys/auth"
	"github.com/labstack/echo/v4"
)

type Server struct {
	router             *echo.Echo
	webauthnController auth.WebAuthnController
	passwordController auth.PasswordController
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
	s.router.StaticFS("/", distDirFS)

	s.router.FileFS("/", "index.html", distIndexHTML)
	s.router.FileFS("/sign-up", "index.html", distIndexHTML)
	s.router.FileFS("/home", "index.html", distIndexHTML)

	s.router.POST("/register/begin", s.webauthnController.BeginRegistration())
	s.router.POST("/register/finish", s.webauthnController.FinishRegistration())
	s.router.POST("/login/begin", s.webauthnController.BeginLogin())
	s.router.POST("/login/finish", s.webauthnController.FinishLogin())

	s.router.POST("/discoverable_login/begin", s.webauthnController.BeginDiscoverableLogin())
	s.router.POST("/discoverable_login/finish", s.webauthnController.FinishDiscoverableLogin())

	s.router.POST("/register/password", s.passwordController.SignUp())
	s.router.POST("/login/password", s.passwordController.Login())
}
