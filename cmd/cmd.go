package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tesarwijaya/ouroboros/internal/config"
	event_repository "github.com/tesarwijaya/ouroboros/internal/domain/event/repository"
	healthz_service "github.com/tesarwijaya/ouroboros/internal/domain/healthz/service"
	player_repository "github.com/tesarwijaya/ouroboros/internal/domain/player/repository"
	player_service "github.com/tesarwijaya/ouroboros/internal/domain/player/service"
	team_repository "github.com/tesarwijaya/ouroboros/internal/domain/team/repository"
	team_service "github.com/tesarwijaya/ouroboros/internal/domain/team/service"
	"github.com/tesarwijaya/ouroboros/internal/entry-point/rest"
	healthz_controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/healthz"
	player_controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/player"
	team_controller "github.com/tesarwijaya/ouroboros/internal/entry-point/rest/controller/team"
	"github.com/tesarwijaya/ouroboros/internal/resource"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
)

func NewCmd() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "server-start",
				Usage: "start the fcking server!",
				Action: func(*cli.Context) error {
					server := newApp(func(lc fx.Lifecycle, server rest.RestServer, db *sql.DB) {
						lc.Append(fx.Hook{
							OnStart: func(ctx context.Context) error {
								go server.Start()

								return nil
							},
							OnStop: func(ctx context.Context) error {
								fmt.Println("closing db...")

								return db.Close()
							},
						})
					})

					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()

					if err := server.Start(ctx); err != nil {
						panic(err)
					}

					<-server.Done()

					ctxStop, cancelStop := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancelStop()

					if err := server.Stop(ctxStop); err != nil {
						panic(err)
					}

					return nil
				},
			},
		},
	}
}

func newApp(invoker ...interface{}) *fx.App {
	return fx.New(
		fx.Provide(
			rest.NewRestServer,
			config.NewConfig,

			resource.NewSQLConnection,
			resource.NewEventStoreConnection,

			healthz_controller.NewHealthzController,
			healthz_service.NewHealthzService,

			player_controller.NewPlayerController,
			player_service.NewPlayerService,
			player_repository.NewPlayerReposity,

			team_controller.NewTeamController,
			team_service.NewTeamService,
			team_repository.NewTeamReposity,

			event_repository.NewTeamReposity,
		),
		fx.Invoke(invoker...),
	)
}
