package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func SetupFrontendRoutes(r *gin.Engine) {
	r.Static("/static", "./frontend")
	r.LoadHTMLGlob("frontend/html/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})
}
