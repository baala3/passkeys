package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	gormsessions "github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var userStore *UserDB
var sessionStore gormsessions.Store

func main() {

	db, err := gorm.Open(sqlite.Open("db/test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	// TODO: use DB
	userStore = DB()

	// session store
	sessionStore = gormsessions.NewStore(db, true, []byte("secret"))
		
	r := gin.Default()

	r.Use(sessions.Sessions("mysession", sessionStore))
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
