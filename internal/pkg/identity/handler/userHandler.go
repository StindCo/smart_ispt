package handler

import (
	"fmt"
	"net/http"

	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service interfaces.UserService
}

func NewUserHandler(app *gin.RouterGroup, auth *jwt.GinJWTMiddleware, service interfaces.UserService) {
	userHandler := UserHandler{
		Service: service,
	}
	app.GET("", userHandler.List)
	app.GET("/:userId", userHandler.GetUser)
	app.PUT("/:userId", userHandler.UpdatePasswordForUser)

	// Need a Admin authorization
	app.Use(auth.MiddlewareFunc())
	{
		app.GET("/admins", userHandler.GetAdminsUsers)
		app.GET("/developpers", userHandler.GetDeveloppersUsers)

		app.POST("", userHandler.CreateUser)
		app.PUT("/:userId/role/:roleId", userHandler.SetRole)

		app.PUT("/:userId/admins", userHandler.SetAdminPermission)
		app.DELETE("/:userId/admins", userHandler.RemoveAdminPermission)

		app.PUT("/:userId/developpers", userHandler.SetDevelopperPermission)
		app.DELETE("/:userId/developpers", userHandler.RemoveDevelopperPermission)

		app.DELETE("/:userId", userHandler.DeleteUser)
	}

}

func (h UserHandler) CreateUser(c *gin.Context) {
	type UserDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Fullname string `json:"fullname"`
	}
	var userDTO UserDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Service.CreateUser(userDTO.Username, userDTO.Password, userDTO.Fullname)
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
		"ressource": fmt.Sprintf("/users/%v", user.ID),
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

func (h UserHandler) GetAdminsUsers(c *gin.Context) {
	users, err := h.Service.GetUsersWhoAreAdmin()
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

func (h UserHandler) GetDeveloppersUsers(c *gin.Context) {
	users, err := h.Service.GetUsersWhoAreDevelopper()
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

func (h UserHandler) SetAdminPermission(c *gin.Context) {
	userId := c.Param("userId")

	user, err := h.Service.SetAdminPermission(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h UserHandler) SetDevelopperPermission(c *gin.Context) {
	userId := c.Param("userId")

	user, err := h.Service.SetDevelopperPermission(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h UserHandler) RemoveAdminPermission(c *gin.Context) {
	userId := c.Param("userId")

	user, err := h.Service.RemoveAdminPermission(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
	})
}

func (h UserHandler) RemoveDevelopperPermission(c *gin.Context) {
	userId := c.Param("userId")

	user, err := h.Service.RemoveDevelopperPermission(userId)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data":   user,
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
