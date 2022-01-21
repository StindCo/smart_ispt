package main

import (
	"github.com/StindCo/smart_ispt/internal/pkg/identity/handler"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/service"
	"github.com/StindCo/smart_ispt/pkg/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LaunchIdentity(app *gin.RouterGroup, db *gorm.DB) {
	// Init repositories
	userRepository := repository.NewUserGORMRepository(db)
	roleRepository := repository.NewRoleGORMRepository(db)

	// Init all services
	userService := service.NewUserService(*userRepository, *roleRepository)
	roleService := service.NewRoleService(*roleRepository, *userRepository)

	// Init All securities applications
	authMiddleware, err := security.InitAdminSecurityMiddleware(*userRepository)
	if err != nil {
		panic("Erreur lors de l'initialisation de la sécurité pour les utilisateurs")
	}

	app.GET("/test/*action", func(c *gin.Context) {
		c.String(200, c.Param("action"))
	})

	app.POST("/login", authMiddleware.LoginHandler)

	handler.NewUserHandler(app.Group("v1/users"), authMiddleware, userService)
	handler.NewRoleHandler(app.Group("v1/roles"), authMiddleware, roleService)
}
