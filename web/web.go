package web

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kataras/blocks"
)

var views *blocks.Blocks

func home(c *gin.Context) {
	buf := new(bytes.Buffer)
	err := views.ExecuteTemplate(buf, "index", "main", nil)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}

	c.String(200, buf.String())
}

func login(c *gin.Context) {
	c.String(200, "Login")
}

func logs(c *gin.Context) {
	c.String(200, "Logs")
}

func removeAutoCheckin(c *gin.Context) {
	c.String(200, "Remove")
}

func addAutoCheckin(c *gin.Context) {
	c.String(200, "Add")
}

func Serve(host string, port int) error {

	views = blocks.New("./views")
	err := views.Load()
	if err != nil {
		log.Fatal("Error loading views: ", err)
	}

	router := gin.Default()

	router.GET("/", home)
	router.POST("/login", login)
	router.GET("/logs", logs)
	router.POST("/remove", removeAutoCheckin)
	router.POST("/add", addAutoCheckin)

	router.Run(fmt.Sprintf("%s:%d", host, port))
	log.Printf("Listening on %s:%d", host, port)

	return nil
}
