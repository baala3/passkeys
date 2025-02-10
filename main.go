package main

import (
	"log"

	"github.com/baala3/passkey-demo/webauthn"

	"github.com/baala3/passkey-demo/routes"
	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


func main() {
	// gin for web framework server
	r := gin.Default()

	// db
	db, err := gorm.Open(sqlite.Open("db/test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	// session store
	sessionStore := gormsessions.NewStore(db, true, []byte("secret"))
	r.Use(sessions.Sessions("mysession", sessionStore))

	// routes
	routes.SetupFrontendRoutes(r)

	// webauthn config
	webauthnController := webauthn.NewWebAuthnController()
	webauthnRoutes := routes.NewWebAuthnRouteController(webauthnController)
	webauthnRoutes.WebAuthnRoutes(r)

	log.Println("Starting server on port 8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
