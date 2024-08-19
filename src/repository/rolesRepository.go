package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/graphql-iam/agent-manager/src/model"
	"github.com/graphql-iam/agent-manager/src/util"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RolesRepository struct {
	db    *mongo.Database
	cache *cache.Cache
}

func NewRolesRepository(db *mongo.Database, c *cache.Cache) *RolesRepository {
	return &RolesRepository{
		db:    db,
		cache: c,
	}
}

func (r *RolesRepository) GetRoleByName(name string) (model.Role, error) {
	res, found := r.cache.Get(name)
	if found {
		return res.(model.Role), nil
	}

	var result model.Role
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := r.db.Collection("rolesWithPolicies").FindOne(ctx, bson.D{{"name", name}}).Decode(&result)
	if err != nil {
		fmt.Println(err.Error())
		return result, errors.New(fmt.Sprintf("could not find role with name %s", name))
	}

	r.cache.Set(name, result, cache.DefaultExpiration)
	return result, nil
}

func (r *RolesRepository) GetRolesByNames(names []string) ([]model.Role, error) {
	var cacheResult []model.Role
	unresolvedNames := names

	for _, name := range names {
		res, found := r.cache.Get(name)
		if found {
			cacheResult = append(cacheResult, res.(model.Role))
			unresolvedNames = util.FilterArray(unresolvedNames, func(s string) bool {
				return s != name
			})
		}
	}
	if len(unresolvedNames) < 1 {
		return cacheResult, nil
	}

	var queryResult []model.Role
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := r.db.Collection("rolesWithPolicies").Find(ctx, bson.D{{"name", bson.D{{"$in", unresolvedNames}}}})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &queryResult)
	if err != nil {
		return nil, err
	}

	for _, role := range queryResult {
		r.cache.Set(role.Name, role, cache.DefaultExpiration)
	}

	return append(cacheResult, queryResult...), nil
}
