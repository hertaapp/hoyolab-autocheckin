Hoyolab Autocheckin
=====================

Auto-checkin service for Hoyolab games, including Honkai Impact 3rd, Genshin Impact 3rd, and Honkai: Star Rail.



## Usage

This tool needs to be hosted on a web server. 


```
NAME:
   hoyolab-autocheckin - Hoyolab auto checkin

USAGE:
   hoyolab-autocheckin [global options] command [command options] [arguments...]

COMMANDS:
   serve    Start a web server
   cron     Start a job scheduler
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

`hoyolab-autocheckin serve` provides the Web UI for users to manage their sessions.

```
NAME:
   hoyolab-autocheckin serve - Start a web server

USAGE:
   hoyolab-autocheckin serve [command options] [arguments...]

OPTIONS:
   --port value  Port to listen on
   --host value  Host to listen on
   --help, -h    show help
```


`hoyolab-autocheckin cron` is the cronjob to actually do the checkin.

```
NAME:
   hoyolab-autocheckin cron - Start a job scheduler

USAGE:
   hoyolab-autocheckin cron [command options] [arguments...]

OPTIONS:
   --help, -h  show help
```


## Development

Run `go run .` to start the dev process.


### Requirements

The following packages are required:

```
gorm.io/gorm                # ORM
gorm.io/driver/mysql        # ORM
github.com/gin-gonic/gin    # web framework
github.com/imroc/req/v3     # send http request
github.com/go-co-op/gocron  # cron
github.com/kataras/blocks   # templating
github.com/urfave/cli/v2    # cmd line
```



