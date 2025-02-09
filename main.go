package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var userStore *UserDB
var sessionStore gormsessions.Store
var webAuthn *webauthn.WebAuthn


func BeginRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := userStore.GetUser(username)

	if err != nil {
		displayName := strings.Split(username, "@")[0]
		user = NewUser(username, displayName)
		userStore.PutUser(user)
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := webAuthn.BeginRegistration(user)

	if err!=nil{
		log.Printf("error beginning registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
		return
	}

	// store session data as marshaled JSON
	session := sessions.Default(c)
	bytes, err:= json.Marshal(sessionData)

	if err != nil {
		log.Printf("error marshaling session data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session data"})
		return
	}

	session.Set("registration", bytes)

	log.Printf("registration options: %v", options)

	c.JSON(http.StatusOK, options)

}

func main() {
	var err error

	// gin for web framework server
	r := gin.Default()
	
	// webauthn config
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Passkey Demo",
		RPID: "localhost",
		RPOrigins: []string{"localhost"},
	})
	if err != nil {
		panic(err)
	}

	// db
	db, err := gorm.Open(sqlite.Open("db/test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	// TODO: use DB
	userStore = DB()

	// session store
	sessionStore = gormsessions.NewStore(db, true, []byte("secret"))
	r.Use(sessions.Sessions("mysession", sessionStore))

	// routes
	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/register/begin/:username", BeginRegistration)

	fmt.Println("Starting server on port 8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
