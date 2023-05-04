package cron

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/imroc/req/v3"
	"gorm.io/gorm"

	"hoyolabautocheckin/model"
)

type ApiInfo struct {
	actId string
	url   string
}

var apiMap = map[string]ApiInfo{
	"genshin": {
		actId: "e202102251931481",
		url:   "https://sg-hk4e-api.hoyolab.com/event/sol/sign?lang=en-us",
	},
	"honkai3rd": {
		actId: "e202110291205111",
		url:   "https://sg-public-api.hoyolab.com/event/mani/sign?lang=en-us",
	},
	"honkaistarrail": {
		actId: "e202303301540311",
		url:   "https://sg-public-api.hoyolab.com/event/luna/os/sign",
	},
}

func Checkin(game string, userId uint64, ltoken string, mhyuuid string) (string, error) {
	apiInfo, ok := apiMap[game]
	if !ok {
		return "", nil
	}

	payload := map[string]interface{}{
		"act_id": apiInfo.actId,
	}

	client := req.C()
	resp, err := client.R().SetBody(&payload).SetCookies(
		&http.Cookie{
			Name:  "_MHYUUID",
			Value: mhyuuid,
		},
		&http.Cookie{
			Name:  "ltoken",
			Value: ltoken,
		},
		&http.Cookie{
			Name:  "ltuid",
			Value: strconv.FormatUint(userId, 10),
		},
	).Post(apiInfo.url)

	if err != nil {
		return "", err
	}

	msg := resp.String()
	db := model.GetDb()
	db.Create(&model.CheckinLog{
		ID:        userId,
		CreatedAt: time.Now(),
		Game:      game,
		Msg:       msg,
	})

	return msg, nil
}

func UserCheckin(userId uint64) {
	log.Println("UserCheckin started. userId: ", userId)
	db := model.GetDb()

	var (
		result       *gorm.DB
		user         model.User
		enabledGames model.EnabledGames
	)

	result = db.First(&user, userId)
	if result.Error != nil {
		return
	}

	result = db.First(&enabledGames, userId)
	if result.Error != nil {
		return
	}

	if enabledGames.GenshinEnabled {
		Checkin("genshin", user.ID, user.Ltoken, user.Mhyuuid)
		time.Sleep(500 * time.Millisecond)
	}

	if enabledGames.Honkai3rdEnabled {
		Checkin("honkai3rd", user.ID, user.Ltoken, user.Mhyuuid)
		time.Sleep(500 * time.Millisecond)
	}

	if enabledGames.HonkaiStarRailEnabled {
		Checkin("honkaistarrail", user.ID, user.Ltoken, user.Mhyuuid)
		time.Sleep(500 * time.Millisecond)
	}
}

func AllUserCheckin() {
	log.Println("AllUserCheckin started.")
	db := model.GetDb()

	users := make([]model.User, 0)
	result := db.Find(&users)

	if result.Error != nil {
		return
	}

	for _, user := range users {
		UserCheckin(user.ID)
	}

}

func StartCron(crontab string) {
	s := gocron.NewScheduler(time.UTC)
	s.Cron(crontab).Do(AllUserCheckin)

	log.Println("Cron started.")
	s.StartBlocking()
}
