package main

import (
	"os"

	"github.com/ravlio/wallet/balance/cmd/commands"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	log.Logger = log.With().Caller().Logger()
	app := cli.NewApp()
	app.Name = "Balance Service"
	app.Commands = []cli.Command{
		commands.Run,
		commands.Migrate,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("cli error")
	}
}
