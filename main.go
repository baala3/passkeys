package main

import (
	"log"
	"net/http"

	"github.com/baala3/passkey-demo/webauthn"

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
	r.Static("/static", "./frontend")
	r.LoadHTMLGlob("frontend/html/**")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	// webauthn config
	wc := webauthn.NewWebAuthnController()
	r.GET("/register/begin/:username", wc.BeginRegistration)
	r.POST("/register/finish/:username", wc.FinishRegistration)
	r.GET("/login/begin/:username", wc.BeginLogin)
	r.POST("/login/finish/:username", wc.FinishLogin)
	log.Println("Starting server on port 8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
