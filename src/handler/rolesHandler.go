package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-iam/agent-manager/src/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

type RolesHandler struct {
	db              *mongo.Database
	rolesRepository *repository.RolesRepository
}

func NewRolesHandler(db *mongo.Database, rolesRepository *repository.RolesRepository) RolesHandler {
	return RolesHandler{
		db:              db,
		rolesRepository: rolesRepository,
	}
}

func (r *RolesHandler) GetRoleByName(c *gin.Context) {
	roleName := c.Query("role")
	if roleName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	role, err := r.rolesRepository.GetRoleByName(roleName)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RolesHandler) GetRolesByNames(c *gin.Context) {
	roleNamesStr, found := c.GetQuery("roles")
	roleNames := strings.Split(roleNamesStr, ",")
	if !found || roleNamesStr == "" || len(roleNames) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	roles, err := r.rolesRepository.GetRolesByNames(roleNames)
	if err != nil {
		fmt.Printf("Error getting roles from database: %v\n", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if roles == nil || len(roles) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, roles)
}
