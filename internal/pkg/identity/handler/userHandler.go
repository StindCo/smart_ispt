package handler

import (
	"fmt"
	"net/http"

	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service interfaces.UserService
}

func NewUserHandler(app *gin.RouterGroup, service interfaces.UserService) {
	userHandler := UserHandler{
		Service: service,
	}
	app.POST("", userHandler.CreateUser)
	app.GET("", userHandler.List)
	app.GET("/:userId", userHandler.GetUser)
	app.PUT("/:userId", userHandler.UpdatePasswordForUser)
	app.PUT("/:userId/role/:roleId", userHandler.SetRole)
	app.DELETE("/:userId", userHandler.DeleteUser)
}

func (h UserHandler) CreateUser(c *gin.Context) {
	type UserDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var userDTO UserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Service.CreateUser(userDTO.Username, userDTO.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"status":    "success",
		"data":      user,
		"ressource": fmt.Sprintf("/users/%v", user.ID.String()),
	})
}

func (h UserHandler) List(c *gin.Context) {
	users, err := h.Service.List()
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "error for formating data",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   users,
	})
	c.Abort()
}

func (h UserHandler) GetUser(c *gin.Context) {
	userId := c.Param("userId")
	users, err := h.Service.GetUser(userId)
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

func (h UserHandler) UpdatePasswordForUser(c *gin.Context) {
	userId := c.Param("userId")

	type passwordDTO struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	var passwords passwordDTO
	if err := c.ShouldBindJSON(&passwords); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	users, err := h.Service.UpdatePassword(userId, passwords.OldPassword, passwords.NewPassword)
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

func (h UserHandler) SetRole(c *gin.Context) {
	userId := c.Param("userId")
	roleId := c.Param("roleId")

	users, err := h.Service.SetRole(userId, roleId)
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

func (h UserHandler) DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	err := h.Service.Delete(userId)
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
