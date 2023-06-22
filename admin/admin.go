package admin

import (
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"

	"hoyolabautocheckin/model"

	"gorm.io/gorm"
)

type GameInfo struct {
	Game        string
	LastLogTime string
	LastLog     string
}

type EnhancedUser struct {
	model.User
	Games []*GameInfo
}

func ListUsers(c *cli.Context) error {
	var (
		userMap      = make(map[uint64]*EnhancedUser)
		users        []*model.User
		enabledGames []*model.EnabledGames
		checkinLogs  []*model.CheckinLog
		result       *gorm.DB
	)

	db := model.GetDb()

	// query database
	result = db.Find(&users)
	if result.Error != nil {
		// print error message
		log.Fatalln("Error finding users: ", result.Error)
	}

	result = db.Find(&enabledGames)
	if result.Error != nil {
		// print error message
		log.Fatalln("Error finding enabled games: ", result.Error)
	}

	maxTsQuery := db.Model(&model.CheckinLog{}).Select("id, game, MAX(created_at) AS ts").Group("id, game")
	result = db.Table("checkin_logs as c, (?) as a", maxTsQuery).Where("c.id = a.id AND c.game = a.game").Find(&checkinLogs)
	if result.Error != nil {
		// print error message
		log.Fatalln("Error finding checkin logs: ", result.Error)
	}

	// integrate data
	for _, user := range users {
		userMap[user.ID] = &EnhancedUser{
			User: *user,
		}
	}

	for _, enabledGame := range enabledGames {
		if enabledGame.GenshinEnabled {
			userMap[enabledGame.ID].Games = append(userMap[enabledGame.ID].Games, &GameInfo{
				Game:        "genshin",
				LastLogTime: "",
				LastLog:     "",
			})
		}

		if enabledGame.Honkai3rdEnabled {
			userMap[enabledGame.ID].Games = append(userMap[enabledGame.ID].Games, &GameInfo{
				Game:        "honkai3rd",
				LastLogTime: "",
				LastLog:     "",
			})
		}

		if enabledGame.HonkaiStarRailEnabled {
			userMap[enabledGame.ID].Games = append(userMap[enabledGame.ID].Games, &GameInfo{
				Game:        "honkaistarrail",
				LastLogTime: "",
				LastLog:     "",
			})
		}
	}

	for _, checkinLog := range checkinLogs {
		for _, gameInfo := range userMap[checkinLog.ID].Games {
			if gameInfo.Game == checkinLog.Game {
				gameInfo.LastLogTime = checkinLog.CreatedAt.Format(time.RFC3339)
				gameInfo.LastLog = checkinLog.Msg
			}
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Created At", "Game", "Last Log Time", "Last Log"})
	for _, user := range userMap {
		t.AppendRow(table.Row{user.ID, user.CreatedAt, user.Games[0].Game, user.Games[0].LastLogTime, user.Games[0].LastLog})
		for _, gameInfo := range user.Games[1:] {
			t.AppendRow(table.Row{"", "", gameInfo.Game, gameInfo.LastLogTime, gameInfo.LastLog})
		}
	}
	t.Render()

	return nil
}
