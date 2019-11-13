package commands

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/config"
	"github.com/ravlio/wallet/tranfser"
	"gopkg.in/urfave/cli.v1"
)

var Migrate = cli.Command{
	Name:  "migrate",
	Usage: "Run migrations",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Configuration file",
			Value: "config.yml",
		}},

	Subcommands: []cli.Command{migrationCmd("up"), migrationCmd("down")},
}

func migrationCmd(dir string) cli.Command {
	ret := cli.Command{
		Action: func(c *cli.Context) error {
			cfg := &tranfser.Config{}

			// use config file if it is passed
			if c.String("config") == "" {
				return errors.New("empty config name")
			}

			err := config.Load(c.String("config"), cfg)
			if err != nil {
				return err
			}

			conn, err := pgx.Connect(context.Background(), cfg.PostgresURL)

			if err != nil {
				return err
			}

			m, err := tranfser.NewMigrations(conn, "transfer")

			if err != nil {
				return err
			}

			if dir == "up" {
				_, err = m.Up()
			} else {
				err = m.Down()
			}

			if err != nil {
				return err
			}

			return nil
		}}

	if dir == "up" {
		ret.Name = "up"

	} else {
		ret.Name = "down"
	}

	return ret
}
