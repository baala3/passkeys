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
	emailController controller.EmailController
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

	s.router.FileFS("/", "index.html", distIndexHTML, middleware.NoAuth)
	s.router.FileFS("/sign-up", "index.html", distIndexHTML, middleware.NoAuth)
	s.router.FileFS("/home", "index.html", distIndexHTML, middleware.Auth)
	s.router.FileFS("/passkeys", "index.html", distIndexHTML, middleware.Auth)
	s.router.FileFS("/delete_account", "index.html", distIndexHTML, middleware.Auth)
	s.router.FileFS("/edit_email", "index.html", distIndexHTML, middleware.Auth)

	s.router.POST("/register/begin", s.webauthnCredentialController.BeginRegistration(), middleware.ConditionalAuth)
	s.router.POST("/register/finish", s.webauthnCredentialController.FinishRegistration(), middleware.ConditionalAuth)
	s.router.GET("/credentials", s.webauthnCredentialController.GetCredentials(), middleware.Auth)
	s.router.DELETE("/credentials", s.webauthnCredentialController.DeleteCredential(), middleware.Auth)

	s.router.POST("/login/begin", s.webauthnAssertionsController.BeginLogin(), middleware.ConditionalAuth)
	s.router.POST("/login/finish", s.webauthnAssertionsController.FinishLogin(), middleware.ConditionalAuth)
	s.router.POST("/discoverable_login/begin", s.webauthnAssertionsController.BeginDiscoverableLogin(), middleware.NoAuth)
	s.router.POST("/discoverable_login/finish", s.webauthnAssertionsController.FinishDiscoverableLogin(), middleware.NoAuth)

	s.router.POST("/register/password", s.passwordController.SignUp(), middleware.NoAuth)
	s.router.POST("/login/password", s.passwordController.Login(), middleware.NoAuth)
	s.router.POST("/logout", s.passwordController.Logout(), middleware.Auth)
	s.router.DELETE("/delete_account", s.passwordController.DeleteAccount(), middleware.Auth)
	s.router.POST("/change_email", s.emailController.ChangeEmail(), middleware.Auth)
}
