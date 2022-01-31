package main

import (
	"github.com/StindCo/smart_ispt/internal/pkg/identity/handler"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/service"
	"github.com/StindCo/smart_ispt/pkg/applogger"
	"github.com/StindCo/smart_ispt/pkg/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Cette fonction permet de dispatcher les endPoints de la sous-application 'Identity'
// elle prend en paramètre une instance de GIN et une instance de GORM
func DispatchIdentityRouter(app *gin.RouterGroup, db *gorm.DB) {

	appLogger := applogger.NewLogger("Identity")
	// Init repositories
	userRepository := repository.NewUserGORMRepository(db)
	roleRepository := repository.NewRoleGORMRepository(db)

	// Init all services
	userService := service.NewUserService(*userRepository, *roleRepository)
	roleService := service.NewRoleService(*roleRepository, *userRepository)

	// Init All securities applications
	authAdminMiddleware, err := security.InitAdminSecurityMiddleware(*userRepository)
	if err != nil {
		panic("Erreur lors de l'initialisation de la sécurité pour les utilisateurs")
	}

	app.POST("/login", authAdminMiddleware.LoginHandler)

	handler.NewUserHandler(app.Group("users"), authAdminMiddleware, userService, appLogger)
	handler.NewRoleHandler(app.Group("roles"), authAdminMiddleware, roleService, appLogger)
}
