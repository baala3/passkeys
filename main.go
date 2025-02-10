package main

import (
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

	if err := StoreSessionData(c, username, sessionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	sessionData, err := LoadSessionData(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credential, err := webAuthn.FinishRegistration(user, *sessionData, c.Request)
	if err != nil {
		log.Printf("error finishing registration: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish registration: " + err.Error()})
		return
	}

	user.AddCredential(*credential)

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func BeginLogin(c *gin.Context) {
	username := c.Param("username")
	user, err := userStore.GetUser(username)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	options, sessionData, err := webAuthn.BeginLogin(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin login"})
		return
	}

	if err := StoreSessionData(c, username, sessionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, options)
}

func FinishLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := userStore.GetUser(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	sessionData, err := LoadSessionData(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// in an actual implementation we should perform additional
	// checks on the returned 'credential'
	_, err = webAuthn.FinishLogin(user, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish login: " + err.Error()})
		return
	}

	// handle successful login
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
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
	r.Static("/static", "./views")
	r.LoadHTMLGlob("views/html/**")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	r.GET("/register/begin/:username", BeginRegistration)
	r.POST("/register/finish/:username", FinishRegistration)
	r.GET("/login/begin/:username", BeginLogin)
	r.POST("/login/finish/:username", FinishLogin)
	log.Println("Starting server on port 8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
