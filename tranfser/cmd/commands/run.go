package commands

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/pkg/config"
	"github.com/ravlio/wallet/pkg/event_bus"
	"github.com/ravlio/wallet/pkg/service"
	"github.com/ravlio/wallet/tranfser"
	"gopkg.in/urfave/cli.v1"
)

var Run = cli.Command{
	Name:  "run",
	Usage: "Run service",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Configuration file",
			Value: "config.yml",
		}},
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

		repo := tranfser.NewRepository(conn)
		// TODO nats
		eb := event_bus.NewLocalBroker()
		// create new service
		svc := tranfser.NewService(repo, eb)

		// create lifecycle
		lc := service.NewLifecycle(
			service.Service(svc),
			service.Name("transfer"),
			service.MetricsPort(cfg.MetricsPort),
		)

		// run service
		return lc.Run()
	},
}
