package main

import (
	"net/http"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/travisbale/birdies-up/internal/booking"
	"github.com/travisbale/birdies-up/internal/clubhouse"
	cli "github.com/urfave/cli/v2"
)

func main() {
	var config booking.Config
	var baseUrl string

	app := &cli.App{
		Name:  "teetimer",
		Usage: "Automate tee time reservations for clubhouselineline-e3.net golf courses",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url",
				Usage:       "Tee sheet URL",
				EnvVars:     []string{"CLUBHOUSE_URL"},
				Destination: &baseUrl,
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Usage:       "Username used to log into the application",
				EnvVars:     []string{"CLUBHOUSE_USERNAME"},
				Destination: &config.Username,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "Password used to log into the application",
				EnvVars:     []string{"CLUBHOUSE_PASSWORD"},
				Destination: &config.Password,
			},
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Configuration file",
				Value:       "config.json",
				EnvVars:     []string{"CLUBHOUSE_CONFIG_FILENAME"},
				Destination: &config.ConfigFileName,
			},
		},
		Action: func(*cli.Context) error {
			config.Client = clubhouse.NewApiClient(baseUrl, http.DefaultClient, log15.New(log15.Ctx{"module": "API client"}))
			config.Logger = log15.New(log15.Ctx{"module": "booking service"})

			service := booking.NewService(config)

			return service.Run()
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err.Error())
	}
}
