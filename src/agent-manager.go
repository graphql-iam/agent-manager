package main

import (
	"github.com/graphql-iam/agent-manager/src/cache"
	"github.com/graphql-iam/agent-manager/src/config"
	"github.com/graphql-iam/agent-manager/src/database"
	"github.com/graphql-iam/agent-manager/src/modules"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.NewConfig),
		fx.Provide(database.NewDatabase),
		fx.Provide(cache.NewCache),
		modules.Repository,
		modules.Handler,
		modules.Server,
	).Run()
}
