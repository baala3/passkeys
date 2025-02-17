package main

import (
	"embed"

	"github.com/baala3/passkey-demo/auth"
	"github.com/labstack/echo/v4"
)

type Server struct {
	router *echo.Echo
	webauthnController auth.WebAuthnController
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

	distDirFS = echo.MustSubFS(dist, "dist")
	distIndexHTML = echo.MustSubFS(indexHTML, "dist")
)

func (s *Server) registerEndpoints() {
	s.router.StaticFS("/", distDirFS)

	s.router.FileFS("/", "index.html", distIndexHTML)
	s.router.FileFS("/sign-up", "index.html", distIndexHTML)

	s.router.GET("/register/begin/:username", s.webauthnController.BeginRegistration())
	s.router.POST("/register/finish/:username", s.webauthnController.FinishRegistration())
	s.router.GET("/login/begin/:username", s.webauthnController.BeginLogin())
	s.router.POST("/login/finish/:username", s.webauthnController.FinishLogin())
}
