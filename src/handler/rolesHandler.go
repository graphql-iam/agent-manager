package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/graphql-iam/agent-manager/src/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

type RolesHandler struct {
	DB              *mongo.Database
	RolesRepository *repository.RolesRepository
}

func (r *RolesHandler) GetRoleByName(c *gin.Context) {
	roleName := c.Query("role")
	if roleName == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	role, err := r.RolesRepository.GetRoleByName(roleName)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, role)
}

func (r *RolesHandler) GetRolesByNames(c *gin.Context) {
	roleNamesStr := c.Query("roles")
	roleNames := strings.Split(roleNamesStr, ",")
	if roleNamesStr == "" || len(roleNames) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	roles, err := r.RolesRepository.GetRolesByNames(roleNames)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if roles == nil || len(roles) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, roles)
}
