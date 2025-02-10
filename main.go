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
	err = session.Save()
	if err != nil {
		log.Printf("error saving session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, options)
}

func FinishRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := userStore.GetUser(username)

	if err != nil {
		log.Printf("error getting user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	session := sessions.Default(c)
	sessionData := session.Get("registration")
	if sessionData == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No registration session found. Please start registration first",
		})
		return
	}

	bytes, ok := sessionData.([]byte)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid session data format",
		})
		return
	}

	var webAuthnSessionData webauthn.SessionData
	if err := json.Unmarshal(bytes, &webAuthnSessionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse session data",
		})
		return
	}

	credential, err := webAuthn.FinishRegistration(user, webAuthnSessionData, c.Request)
	if err != nil {
		log.Printf("error finishing registration: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish registration"})
		return
	}

	user.AddCredential(*credential)

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func main() {
	var err error

	// gin for web framework server
	r := gin.Default()
	
	// webauthn config
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Passkey Demo",
		RPID: "localhost",
		RPOrigins: []string{"http://localhost:8080"},
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
	r.POST("/register/finish/:username", FinishRegistration)
	fmt.Println("Starting server on port 8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
