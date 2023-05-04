package model

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func Connect() *gorm.DB {
	// connect to database
	db, err := gorm.Open(sqlite.Open("db/hoyolab-autocheckin.sqlite3"), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}

func GetDb() *gorm.DB {
	if db == nil {
		db = Connect()
	}
	return db
}
