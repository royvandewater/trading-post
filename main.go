package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/royvandewater/trading-post/auth0creds"
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
			Name:   "auth0-callback-url",
			EnvVar: "AUTH0_CALLBACK_URL",
			Usage:  "Where to have Auth0 redirect after an auth attempt",
		},
		cli.StringFlag{
			Name:   "auth0-client-id",
			EnvVar: "AUTH0_CLIENT_ID",
			Usage:  "Auth0 client id (from auth0.com)",
		},
		cli.StringFlag{
			Name:   "auth0-client-secret",
			EnvVar: "AUTH0_CLIENT_SECRET",
			Usage:  "Auth0 client secret (from auth0.com)",
		},
		cli.StringFlag{
			Name:   "auth0-domain",
			EnvVar: "AUTH0_DOMAIN",
			Usage:  "Auth0 domain (from auth0.com)",
		},
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
	auth0CallbackURL, auth0ClientID, auth0ClientSecret, auth0Domain, mongoDBURL, port := getOpts(context)
	auth0Creds := auth0creds.Auth0Creds{
		CallbackURL:  auth0CallbackURL,
		ClientID:     auth0ClientID,
		ClientSecret: auth0ClientSecret,
		Domain:       auth0Domain,
	}

	server := tradingpostserver.New(auth0Creds, mongoDBURL, port)
	fmt.Printf("Listening on 0.0.0.0:%v\n", port)
	err := server.Run()
	log.Fatalln(err)
}

func getOpts(context *cli.Context) (string, string, string, string, string, int) {
	auth0CallbackURL := context.String("auth0-callback-url")
	auth0ClientID := context.String("auth0-client-id")
	auth0ClientSecret := context.String("auth0-client-secret")
	auth0Domain := context.String("auth0-domain")
	mongoDBURL := context.String("mongodb-url")
	port := context.Int("port")

	if port == 0 || mongoDBURL == "" || auth0CallbackURL == "" || auth0ClientID == "" || auth0ClientSecret == "" || auth0Domain == "" {
		cli.ShowAppHelp(context)

		if auth0CallbackURL == "" {
			color.Red("  Missing required flag --auth0-callback-url or AUTH0_CALLBACK_URL")
		}
		if auth0ClientID == "" {
			color.Red("  Missing required flag --auth0-client-id or AUTH0_CLIENT_ID")
		}
		if auth0ClientSecret == "" {
			color.Red("  Missing required flag --auth0-client-secret or AUTH0_CLIENT_SECRET")
		}
		if auth0Domain == "" {
			color.Red("  Missing required flag --auth0-domain or AUTH0_DOMAIN")
		}
		if mongoDBURL == "" {
			color.Red("  Missing required flag --mongodb-url or MONGODB_URL")
		}
		if port == 0 {
			color.Red("  Missing required flag --port or PORT")
		}
		os.Exit(1)
	}

	return auth0CallbackURL, auth0ClientID, auth0ClientSecret, auth0Domain, mongoDBURL, port
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
