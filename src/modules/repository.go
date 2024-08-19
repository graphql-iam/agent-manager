package modules

import (
	"github.com/graphql-iam/agent-manager/src/repository"
	"go.uber.org/fx"
)

var Repository = fx.Module("repository",
	fx.Provide(repository.NewRolesRepository),
)
