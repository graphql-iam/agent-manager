package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-iam/agent-manager/src/config"
	"github.com/graphql-iam/agent-manager/src/handler"
	"github.com/graphql-iam/agent-manager/src/repository"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

const ConfigPathEnvName = "AGENT_MANAGER_CONFIG_PATH"

func main() {
	configPath := "./config.yaml"

	if os.Getenv(ConfigPathEnvName) != "" {
		configPath = os.Getenv(ConfigPathEnvName)
	}

	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoUrl))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database("graphql-iam")

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

	err = db.CreateView(context.TODO(), "rolesWithPolicies", "roles", pipeline)
	if err != nil {
		log.Fatal(err)
	}

	expire := time.Duration(cfg.CacheOptions.Expiration) * time.Minute
	purge := time.Duration(cfg.CacheOptions.Purge) * time.Minute
	c := cache.New(expire, purge)

	rolesRepository := repository.RolesRepository{
		DB:    db,
		Cache: c,
	}

	rolesHandler := handler.RolesHandler{
		DB:              db,
		RolesRepository: &rolesRepository,
	}

	r := gin.Default()
	r.GET("/role", rolesHandler.GetRoleByName)
	r.GET("/roles", rolesHandler.GetRolesByNames)
	r.GET("/ping", handler.Ping)

	err = r.Run(fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		panic(err)
	}
}
