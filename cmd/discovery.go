package main

import (
	"github.com/StindCo/smart_ispt/internal/pkg/discovery/handler"
	"github.com/StindCo/smart_ispt/internal/pkg/discovery/repository"
	"github.com/StindCo/smart_ispt/internal/pkg/discovery/service"
	identityRepo "github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
	"github.com/StindCo/smart_ispt/pkg/applogger"
	"github.com/StindCo/smart_ispt/pkg/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DispatcheDiscoveryRouter(app *gin.RouterGroup, db *gorm.DB) {

	appLogger := applogger.NewLogger("discovery")

	// Init repositories
	applicationRepo := repository.NewApplicationGORMRepository(db)
	userRepository := identityRepo.NewUserGORMRepository(db)
	roleRepository := identityRepo.NewRoleGORMRepository(db)

	applicationService := service.NewApplicationService(*applicationRepo, *userRepository, *roleRepository)

	// Init All securities applications
	authDevelopperMiddleware, err := security.InitDevelopperSecurityMiddleware(*userRepository)
	if err != nil {
		panic("Erreur lors de l'initialisation de la sécurité pour les utilisateurs")
	}
	handler.NewApplicationHandler(app.Group("applications"), authDevelopperMiddleware, applicationService, appLogger)
}
