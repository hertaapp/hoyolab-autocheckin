package web

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kataras/blocks"
)

var (
	//go:embed views/*
	embedFs embed.FS
	views   *blocks.Blocks
)

func home(c *gin.Context) {
	buf := new(bytes.Buffer)
	err := views.ExecuteTemplate(buf, "index", "main", nil)
	if err != nil {
		log.Fatal("Error executing template: ", err)
	}

	c.Data(200, "text/html", buf.Bytes())
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

	viewsDir, err := fs.Sub(embedFs, "views")
	if err != nil {
		log.Fatal("Error getting views directory: ", err)
	}

	views = blocks.New(viewsDir)
	err = views.Load()

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
