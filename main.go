package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var userStore *UserDB

func main() {

	userStore = DB()
		
	r := gin.Default()
	fmt.Println("Starting server on port 8080")

	// routes
	r.LoadHTMLGlob("views/*")
	r.GET("/", func(c *gin.Context) {
		username := c.Query("username")
		userStore.PutUser(NewUser(username, username, []string{"testcredential"}))
		fmt.Println(userStore.GetUserCount())
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
