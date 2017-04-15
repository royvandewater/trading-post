package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/royvandewater/trading-post/tradingpostserver"
	"github.com/urfave/cli"
	De "github.com/visionmedia/go-debug"
)

var debug = De.Debug("trading-post:main")

func main() {
	app := cli.NewApp()
	app.Name = "trading-post"
	app.Version = version()
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mongodb-url, m",
			EnvVar: "MONGODB_URL",
			Usage:  "Mongo db url to use for data persistence",
			Value:  "mongodb://localhost:27017",
		},
		cli.IntFlag{
			Name:   "port, p",
			EnvVar: "PORT",
			Usage:  "Port to listen on",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	port, mongoDBURL := getOpts(context)
	server := tradingpostserver.New(port, mongoDBURL)
	fmt.Printf("Listening on 0.0.0.0:%v\n", port)
	err := server.Run()
	log.Fatalln(err)
}

func getOpts(context *cli.Context) (int, string) {
	mongoDBURL := context.String("mongodb-url")
	port := context.Int("port")

	if port == 0 || mongoDBURL == "" {
		cli.ShowAppHelp(context)

		if mongoDBURL == "" {
			color.Red("  Missing required flag --mongodb-url or MONGODB_URL")
		}
		if port == 0 {
			color.Red("  Missing required flag --port or PORT")
		}
		os.Exit(1)
	}

	return port, mongoDBURL
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
