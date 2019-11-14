package commands

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/ravlio/wallet/account"
	"github.com/ravlio/wallet/balance"
	"github.com/ravlio/wallet/pkg/config"
	"github.com/ravlio/wallet/pkg/event_bus"
	"github.com/ravlio/wallet/pkg/service"
	"gopkg.in/urfave/cli.v1"
)

var Run = cli.Command{
	Name:  "run",
	Usage: "Run service",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Configuration file",
			Value: "config.yaml",
		}},
	Action: func(c *cli.Context) error {
		cfg := &balance.Config{}

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

		repo := balance.NewRepository(conn)
		accrepo := account.NewRepository(conn)
		// TODO grpc
		acc := account.NewService(accrepo)
		// TODO nats
		eb := event_bus.NewLocalBroker()
		// create new service
		svc, err := balance.NewService(repo, acc, eb)

		if err != nil {
			return err
		}

		// create lifecycle
		lc := service.NewLifecycle(
			service.Service(svc),
			service.Name("balance"),
			service.MetricsPort(cfg.MetricsPort),
		)

		// run service
		return lc.Run()
	},
}
