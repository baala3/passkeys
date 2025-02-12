package main

import (
	"github.com/baala3/passkey-demo/auth"
	"github.com/labstack/echo/v4"
)

type Server struct {
	router *echo.Echo
	webauthnController auth.WebAuthnController
}

func (s *Server) Start() {
	s.registerEndpoints()
	s.router.Logger.Fatal(s.router.Start(":8080"))
}

func (s *Server) registerEndpoints() {
	s.router.Static("/static", "web")
	s.router.File("/", "web/login.html")
	s.router.File("/home", "web/home.html")

	s.router.GET("/register/begin/:username", s.webauthnController.BeginRegistration())
	s.router.POST("/register/finish/:username", s.webauthnController.FinishRegistration())
	s.router.GET("/login/begin/:username", s.webauthnController.BeginLogin())
	s.router.POST("/login/finish/:username", s.webauthnController.FinishLogin())
}
