package handler

import (
	"net/http"

	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	Service interfaces.RoleService
}

func NewRoleHandler(app *gin.RouterGroup, service interfaces.RoleService) {
	roleHandler := RoleHandler{
		Service: service,
	}
	app.POST("", roleHandler.CreateRole)
	app.GET("", roleHandler.List)
	app.GET("/:roleId", roleHandler.GetRole)
	app.GET("/:roleId/users", roleHandler.GetUsersForRole)
	app.PUT("/:roleId", roleHandler.DeleteRole)
}

func (h RoleHandler) CreateRole(c *gin.Context) {
	type RoleDTO struct {
		Name        string `json:"name"`
		Tag         string `json:"tag"`
		Description string `json:"description"`
	}
	var roleDTO RoleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role, err := h.Service.CreateRole(roleDTO.Name, roleDTO.Tag, roleDTO.Description)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   role,
	})
}

func (h RoleHandler) List(c *gin.Context) {
	roles, err := h.Service.List()
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   roles,
	})
	c.Abort()
}

func (h RoleHandler) GetRole(c *gin.Context) {
	roleID := c.Param("roleId")
	role, err := h.Service.GetRole(roleID)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   role,
	})
}

func (h RoleHandler) GetUsersForRole(c *gin.Context) {
	roleID := c.Param("roleId")
	users, err := h.Service.GetUsers(roleID)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   users,
	})
}

func (h RoleHandler) DeleteRole(c *gin.Context) {
	roleId := c.Param("roleId")
	err := h.Service.Delete(roleId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
	})
}
