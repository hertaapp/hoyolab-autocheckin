package web

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"strconv"
	"strings"
	"time"

	"hoyolabautocheckin/model"

	"github.com/gin-gonic/gin"
	"github.com/kataras/blocks"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var (
	//go:embed views/*
	embedFs embed.FS
	views   *blocks.Blocks
)

func render(c *gin.Context, name string, data interface{}) {
	buf := new(bytes.Buffer)
	err := views.ExecuteTemplate(buf, name, "main", data)
	if err != nil {
		log.Print("Error executing template: ", err)
	}

	c.Data(200, "text/html", buf.Bytes())
}

func requireCookies(c *gin.Context) *model.User {
	ltuid, err := c.Cookie("ltuid")
	if err != nil || ltuid == "" {
		return nil
	}

	userId, err := strconv.ParseUint(ltuid, 10, 64)
	if err != nil {
		return nil
	}

	ltoken, err := c.Cookie("ltoken")
	if err != nil || ltoken == "" {
		return nil
	}

	mhyuuid, err := c.Cookie("_MHYUUID")
	if err != nil || mhyuuid == "" {
		return nil
	}

	return &model.User{
		ID:        userId,
		Ltoken:    ltoken,
		Mhyuuid:   mhyuuid,
		CreatedAt: time.Now(),
	}
}

func home(c *gin.Context) {
	render(c, "index", nil)
}

func login(c *gin.Context) {
	hoyolabCookie, exist := c.GetPostForm("hoyolab-cookie")
	if !exist {
		c.Redirect(302, "/")
	}

	var user model.User

	parts := strings.Split(hoyolabCookie, ";")
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "_MHYUUID":
			user.Mhyuuid = kv[1]
		case "ltoken":
			user.Ltoken = kv[1]
		case "ltuid":
			userId, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				user.ID = 0
				continue
			}
			user.ID = userId
		}
	}

	// invalid data, return to home page
	if user.Mhyuuid == "" || user.Ltoken == "" || user.ID == 0 {
		c.Redirect(302, "/")
	}

	// store the data temporarily in cookie so that we can use it later
	c.SetCookie("ltuid", strconv.FormatUint(user.ID, 10), 0, "", "", false, true)
	c.SetCookie("ltoken", user.Ltoken, 0, "", "", false, true)
	c.SetCookie("_MHYUUID", user.Mhyuuid, 0, "", "", false, true)

	c.Redirect(302, "/logs")
}

func logs(c *gin.Context) {
	cookieUser := requireCookies(c)
	if cookieUser == nil {
		c.Redirect(302, "/")
		return
	}

	// find user in database
	db := model.GetDb()

	var (
		result       *gorm.DB
		user         model.User
		enabledGames model.EnabledGames
		checkinLogs  []model.CheckinLog
	)

	result = db.First(&user, cookieUser.ID)

	notFound := false
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			notFound = true
		} else {
			log.Println("Error finding user: ", result.Error)
			c.AbortWithStatus(500)
			return
		}
	}

	if notFound {
		// if not found, set default values to empty
		enabledGames.ID = cookieUser.ID
		enabledGames.GenshinEnabled = false
		enabledGames.Honkai3rdEnabled = false
		enabledGames.HonkaiStarRailEnabled = false

		checkinLogs = make([]model.CheckinLog, 0)

	} else {
		// if user found, find enabled games and checkin logs
		result = db.First(&enabledGames, cookieUser.ID)
		if result.Error != nil {
			enabledGames.ID = user.ID

			enabledGames.GenshinEnabled = false
			enabledGames.Honkai3rdEnabled = false
			enabledGames.HonkaiStarRailEnabled = false
		}

		result = db.Where("id = ?", cookieUser.ID).Order("created_at desc").Find(&checkinLogs)
		if result.Error != nil {
			checkinLogs = make([]model.CheckinLog, 0)
		}
	}

	// handle message from last page
	msg, err := c.Cookie("msg")
	if err == nil {
		c.SetCookie("msg", "", -1, "", "", false, true)
	}

	render(c, "logs", map[string]interface{}{
		"msg":          msg,
		"ltuid":        cookieUser.ID,
		"notFound":     notFound,
		"enabledGames": enabledGames,
		"checkinLogs":  checkinLogs,
	})
}

func removeAutoCheckin(c *gin.Context) {
	cookieUser := requireCookies(c)

	db := model.GetDb()

	db.Delete(&model.User{}, cookieUser.ID)
	db.Delete(&model.EnabledGames{}, cookieUser.ID)

	c.SetCookie("msg", "Auto checkin settings deleted", 0, "", "", false, true)

	c.Redirect(302, "/logs")
}

func updateAutoCheckin(c *gin.Context) {
	cookieUser := requireCookies(c)
	if cookieUser == nil {
		c.Redirect(302, "/")
		return
	}

	games := c.PostFormArray("games[]")
	if len(games) == 0 {
		c.Redirect(302, "/logs")
		return
	}

	enabledGames := model.EnabledGames{
		ID:                    cookieUser.ID,
		GenshinEnabled:        slices.Contains(games, "genshin"),
		Honkai3rdEnabled:      slices.Contains(games, "honkai3rd"),
		HonkaiStarRailEnabled: slices.Contains(games, "honkaistarrail"),
	}

	db := model.GetDb()
	db.Save(&cookieUser)
	db.Save(&enabledGames)

	c.SetCookie("msg", "Successfully updated auto checkin settings", 0, "", "", false, true)

	c.Redirect(302, "/logs")
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
	router.POST("/update", updateAutoCheckin)

	router.Run(fmt.Sprintf("%s:%d", host, port))
	log.Printf("Listening on %s:%d", host, port)

	return nil
}
