package auth

import (
	"log"
	"net/http"

	"github.com/baala3/passkey-demo/users"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnController struct {
	WebAuthnAPI *webauthn.WebAuthn
	UserStore users.UserRepository
}

func (wc *WebAuthnController) BeginRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := wc.UserStore.GetUser(username)

	if err != nil {
		wc.UserStore.PutUser(username)
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := wc.WebAuthnAPI.BeginRegistration(user)

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

func (wc *WebAuthnController) FinishRegistration(c *gin.Context) {
	username := c.Param("username")

	user, err := wc.UserStore.GetUser(username)

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

	credential, err := wc.WebAuthnAPI.FinishRegistration(user, *sessionData, c.Request)
	if err != nil {
		log.Printf("error finishing registration: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish registration: " + err.Error()})
		return
	}

	user.AddCredential(*credential)

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func (wc *WebAuthnController) BeginLogin(c *gin.Context) {
	username := c.Param("username")
	user, err := wc.UserStore.GetUser(username)
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user"})
		return
	}

	options, sessionData, err := wc.WebAuthnAPI.BeginLogin(user)
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

func (wc *WebAuthnController) FinishLogin(c *gin.Context) {
	username := c.Param("username")

	user, err := wc.UserStore.GetUser(username)
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
	_, err = wc.WebAuthnAPI.FinishLogin(user, *sessionData, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to finish login: " + err.Error()})
		return
	}

	// handle successful login
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}
