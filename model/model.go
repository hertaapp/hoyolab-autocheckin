package model

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"hoyolabautocheckin/utils"
)

// _MHYUUID=1bea8564-1666-4776-8d05-42d27332a66f;ltoken=2Fo4TM35AkLyzkdJDXmGWnaom4rYcD1gGfjyX43a;ltuid=15928933
type User struct {
	ID        uint64 `gorm:"primarykey"`
	Mhyuuid   string
	Ltoken    string
	CreatedAt time.Time
}

type EnabledGames struct {
	ID                    uint64 `gorm:"primarykey"`
	GenshinEnabled        bool
	Honkai3rdEnabled      bool
	HonkaiStarRailEnabled bool
}

type CheckinLog struct {
	ID        uint64    `gorm:"primarykey"`
	Game      string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"primaryKey"`
	Msg       string
}

var db *gorm.DB

func Connect(host string, port int, username string, password string, database string, charset string, parseTime bool, loc string) *gorm.DB {
	// connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}

func GetDb() *gorm.DB {
	if db == nil {
		host := utils.Getenv("MYSQL_HOST", "localhost")
		port := utils.GetenvInt("MYSQL_PORT", 3306)
		username := utils.Getenv("MYSQL_USERNAME", "root")
		password := utils.Getenv("MYSQL_PASSWORD", "")
		database := utils.Getenv("MYSQL_DATABASE", "hoyolab_autocheckin")

		db = Connect(host, port, username, password, database, "utf8mb4", true, "Local")
	}
	return db
}
