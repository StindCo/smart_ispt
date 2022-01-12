package handler

import (
	"net/http"

	"github.com/StindCo/smart_ispt/internal/pkg/identity/interfaces"
	"github.com/gin-gonic/gin"
)

func NewUserHandler(app *gin.RouterGroup, service interfaces.UserService) {
	app.PUT("/", CreateUser(service))
}

func CreateUser(service interfaces.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		type UserDTO struct {
			username string `json:"username"`
			password string `json:"password"`
		}
		var userDTO UserDTO
		if err := c.ShouldBindJSON(&userDTO); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// user, err := service.CreateUser(userDTO.username, userDTO.password)
		// if err != nil {
		// 	c.JSON(400, gin.H{
		// 		"status":  "error",
		// 		"message": "error for formating data",
		// 	})
		// 	c.Abort()
		// 	return
		// }
		// c.JSON(200, gin.H{
		// 	"status": "success",
		// 	"data":   user,
		// })
		// c.Abort()
	}
}
