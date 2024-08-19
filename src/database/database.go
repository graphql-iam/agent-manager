package database

import (
	"context"
	"github.com/graphql-iam/agent-manager/src/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

func NewDatabase(lc fx.Lifecycle, cfg config.Config) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoUrl))
	if err != nil {
		panic(err)
	}
	db := client.Database("graphql-iam")

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return prepareDb(db, ctx)
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return db
}

func prepareDb(db *mongo.Database, ctx context.Context) error {
	pipeline := mongo.Pipeline{
		{
			{"$lookup", bson.D{
				{"from", "policies"},
				{"localField", "policyIds"},
				{"foreignField", "id"},
				{"as", "policies"},
			}},
		},
		{
			{"$project", bson.D{
				{"name", 1},
				{"policies", 1},
				{"_id", 0},
			}},
		},
	}

	return db.CreateView(context.TODO(), "rolesWithPolicies", "roles", pipeline)
}
