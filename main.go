package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"hoyolabautocheckin/cron"
	"hoyolabautocheckin/model"
	"hoyolabautocheckin/web"
)

func serve(c *cli.Context) error {
	// start a web server on the specified host and port
	web.Serve(c.String("host"), c.Int("port"))

	return nil
}

func runCron(c *cli.Context) error {
	cron.StartCron()
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
		},
	}

	// setup db
	db := model.GetDb()
	db.AutoMigrate(&model.User{}, &model.EnabledGames{}, &model.CheckinLog{})
	log.Printf("Database connected: %v", db)

	app.Run(os.Args)

}
