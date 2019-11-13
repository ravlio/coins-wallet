package commands

import (
	"gopkg.in/urfave/cli.v1"
)

var Migrate = cli.Command{
	Name:        "migrate",
	Description: "Run migrations",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Configuration file",
			Value: "config.yml",
		}},
	Action: func(c *cli.Context) error {
		// TODO migrations
	}}
