package routes

import (
	"github.com/baala3/passkey-demo/webauthn"
	"github.com/gin-gonic/gin"
)
type WebAuthnRouteController struct {
	webAuthnController webauthn.WebAuthnController
}
func NewWebAuthnRouteController(webauthnController webauthn.WebAuthnController) WebAuthnRouteController {
	return WebAuthnRouteController{webauthnController}
}
func (rc *WebAuthnRouteController) WebAuthnRoutes(r *gin.Engine) {
	r.GET("/register/begin/:username", rc.webAuthnController.BeginRegistration)
	r.POST("/register/finish/:username",rc.webAuthnController.FinishRegistration)
	r.GET("/login/begin/:username", rc.webAuthnController.BeginLogin)
	r.POST("/login/finish/:username", rc.webAuthnController.FinishLogin)
}