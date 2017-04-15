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
		cli.IntFlag{
			Name:   "port, p",
			EnvVar: "PORT",
			Usage:  "Port to listen on",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	port := getOpts(context)
	server := tradingpostserver.New(port)
	fmt.Printf("Listening on 0.0.0.0:%v\n", port)
	err := server.Run()
	log.Fatalln(err)
}

func getOpts(context *cli.Context) int {
	port := context.Int("port")

	if port == 0 {
		cli.ShowAppHelp(context)

		if port == 0 {
			color.Red("  Missing required flag --port or PORT")
		}
		os.Exit(1)
	}

	return port
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
