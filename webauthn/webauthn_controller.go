package webauthn

import (
	"log"
	"net/http"

	"github.com/baala3/passkey-demo/users"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnController interface {
	BeginRegistration(c *gin.Context)
	FinishRegistration(c *gin.Context)
	BeginLogin(c *gin.Context)
	FinishLogin(c *gin.Context)
}

type webAuthnController struct {
	webAuthn *webauthn.WebAuthn
	userStore users.UserRepository
}

func NewWebAuthnController() *webAuthnController {
	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "Passkey Demo",
		RPID: "localhost",
		RPOrigins: []string{"http://localhost:8080"},
	})
	if err != nil {
		log.Fatalf("error creating webauthn: %v", err)
	}

	return &webAuthnController{
		webAuthn: webAuthn,
		userStore: users.NewUserRepository(), // TODO: use DB
	}
}

func (wc *webAuthnController) BeginRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := wc.userStore.GetUser(username)

	if err != nil {
		wc.userStore.PutUser(username)
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := wc.webAuthn.BeginRegistration(user)

	if err!=nil{
		log.Printf("error beginning registration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin registration"})
		return
	}

	if err := storeSessionData(c, username, sessionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, options)
}

func (wc *webAuthnController) FinishRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := wc.userStore.GetUser(username)

	if err != nil {
		log.Printf("error getting user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	sessionData, err := loadSessionData(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credential, err := wc.webAuthn.FinishRegistration(user, *sessionData, c.Request)
	if err != nil {
		log.Printf("error finishing registration: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish registration: " + err.Error()})
		return
	}

	user.AddCredential(*credential)

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func (wc *webAuthnController) BeginLogin(c *gin.Context) {
	username := c.Param("username")
	user, err := wc.userStore.GetUser(username)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	options, sessionData, err := wc.webAuthn.BeginLogin(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin login"})
		return
	}

	if err := storeSessionData(c, username, sessionData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, options)
}

func (wc *webAuthnController) FinishLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := wc.userStore.GetUser(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	sessionData, err := loadSessionData(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// in an actual implementation we should perform additional
	// checks on the returned 'credential'
	_, err = wc.webAuthn.FinishLogin(user, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish login: " + err.Error()})
		return
	}

	// handle successful login
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
