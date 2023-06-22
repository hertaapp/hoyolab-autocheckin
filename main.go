package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"hoyolabautocheckin/admin"
	"hoyolabautocheckin/cron"
	"hoyolabautocheckin/model"
	"hoyolabautocheckin/web"
)

func serve(c *cli.Context) error {
	// setup db
	db := model.GetDb()
	db.AutoMigrate(&model.User{}, &model.EnabledGames{}, &model.CheckinLog{})
	log.Printf("Database connected")

	// start a web server on the specified host and port
	web.Serve(c.String("host"), c.Int("port"))

	return nil
}

func runCron(c *cli.Context) error {
	// setup db
	db := model.GetDb()
	db.AutoMigrate(&model.User{}, &model.EnabledGames{}, &model.CheckinLog{})
	log.Printf("Database connected")

	crontab := os.Getenv("CRONTAB")
	if crontab == "" {
		crontab = "5 17 * * *"
	}

	cron.StartCron(crontab)
	return nil
}

func main() {
	// create a CLI app with github.com/urfave/cli
	// the app should support two sub commands:
	// - `serve` which starts a web server
	// - `cron` which starts a job scheduler

	// the `serve` command should support the following flags:
	// - `--port` which specifies the port to listen on
	// - `--host` which specifies the host to listen on

	app := &cli.App{
		Name:  "hoyolab-autocheckin",
		Usage: "Hoyolab auto checkin",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start a web server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "8080",
						Usage: "Port to listen on",
					},
					&cli.StringFlag{
						Name:  "host",
						Value: "localhost",
						Usage: "Host to listen on",
					},
				},
				Action: serve,
			},

			{
				Name:   "cron",
				Usage:  "Start a job scheduler",
				Action: runCron,
			},

			{
				Name:    "admin",
				Aliases: []string{"a"},
				Usage:   "Run admin tool",
				Subcommands: []*cli.Command{
					{
						Name:   "ls",
						Usage:  "List all users",
						Action: admin.ListUsers,
					},

					{
						Name:  "checkin",
						Usage: "Checkin for specified user. Used for test.",
						Flags: []cli.Flag{
							&cli.Uint64Flag{
								Name:     "uid",
								Aliases:  []string{"u"},
								Usage:    "User ID",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							uid := c.Uint64("uid")
							log.Printf("Checking in for user %s", uid)
							cron.UserCheckin(uid)
							return nil
						},
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
