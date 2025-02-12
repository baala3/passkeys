package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/baala3/passkey-demo/auth"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	webauthnController auth.WebAuthnController
}

func (s *Server) Start() {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("./db/%s.db", os.Getenv("DB_NAME"))), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//session store
	sessionStore := gormsessions.NewStore(db, true, []byte(os.Getenv("SESSION_SECRET")))
	s.router.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), sessionStore))

	//routes
	s.registerEndpoints()
	log.Println("Starting server on port 8080")
	_ = s.router.Run(":8080")
}

func (s *Server) registerEndpoints() {
	s.router.Static("/static", "web")
	s.router.LoadHTMLGlob("web/*")
	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	s.router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	s.router.GET("/register/begin/:username", s.webauthnController.BeginRegistration)
	s.router.POST("/register/finish/:username", s.webauthnController.FinishRegistration)
	s.router.GET("/login/begin/:username", s.webauthnController.BeginLogin)
	s.router.POST("/login/finish/:username", s.webauthnController.FinishLogin)
}
