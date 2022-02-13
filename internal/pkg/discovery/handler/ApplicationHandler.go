package handler

import (
	"fmt"
	"net/http"

	"github.com/StindCo/smart_ispt/internal/entities"
	"github.com/StindCo/smart_ispt/internal/pkg/discovery/interfaces"
	"github.com/StindCo/smart_ispt/pkg/id"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
)

type ApplicationHandler struct {
	Service interfaces.ApplicationService
	Logger  hclog.Logger
}

func getUserActor(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	userId := claims["id"].(string)
	return userId
}

func NewApplicationHandler(app *gin.RouterGroup, auth *jwt.GinJWTMiddleware, service interfaces.ApplicationService, logger hclog.Logger) {
	ApplicationHandler := ApplicationHandler{
		Service: service,
		Logger:  logger,
	}

	// TODO: Il faudrait penser à créer des Querys pour les filtrages
	app.GET("", ApplicationHandler.List)

	app.GET("/:applicationID", ApplicationHandler.GetOneApplication)
	app.GET("/:applicationID/developpers", ApplicationHandler.GetApplicationDeveloppers)
	app.GET("/:applicationID/consumers", ApplicationHandler.GetApplicationConsumers)

	app.GET("/developpers/:developperID/applications", ApplicationHandler.GetApplicationsCreatedByDevelopperID)
	app.GET("/consumers/:roleID/applications", ApplicationHandler.GetApplicationsCreatedForRoleID)

	// Need a Admin authorization
	app.Use(auth.MiddlewareFunc())
	{
		app.POST("", ApplicationHandler.CreateApplication)
		app.PUT("/:applicationID/consumers", ApplicationHandler.AddApplicationConsumers)
		app.PUT("/:applicationID/developpers", ApplicationHandler.AddApplicationDevelopper)
		app.DELETE("/:applicationID", ApplicationHandler.DeleteApplication)
	}
}

func (ah ApplicationHandler) CreateApplication(c *gin.Context) {

	userDevelopperID := getUserActor(c)

	type ApplicationDTO struct {
		Name        string `json:"name"`
		PowerBy     string `json:"power_by"`
		SmartName   string `json:"smart_name"`
		DomainName  string `json:"domain_name"`
		TestPath    string `json:"test_path"`
		UrlPath     string `json:"url_path"`
		Ip          string `json:"ip"`
		Description string `json:"description"`
	}
	var applicationDTO ApplicationDTO
	if err := c.ShouldBindJSON(&applicationDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application := &entities.Application{
		ID:          id.NewID().String(),
		Name:        applicationDTO.Name,
		UrlPath:     applicationDTO.UrlPath,
		SmartName:   applicationDTO.SmartName,
		PowerBy:     applicationDTO.PowerBy,
		TestPath:    applicationDTO.TestPath,
		Description: applicationDTO.Description,
		DomainName:  applicationDTO.DomainName,
		Ip:          applicationDTO.Ip,
	}

	applicationCreated, err := ah.Service.CreateApplication(userDevelopperID, application)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// h.Logger.Info(fmt.Sprintf("Admin '%v' create a user with username = %v and id = %v", userActor.Username, user.Username, user.ID))
	//TODO: "à faire"
	c.JSON(201, gin.H{
		"status":    "success",
		"data":      applicationCreated,
		"ressource": fmt.Sprintf("/users/%v", applicationCreated.ID),
	})
}

func (ah ApplicationHandler) GetOneApplication(c *gin.Context) {
	appID := c.Param("applicationID")
	application, err := ah.Service.GetOneApplication(appID)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   application,
	})
}

func (ah ApplicationHandler) GetApplicationDeveloppers(c *gin.Context) {
	appID := c.Param("applicationID")
	developpers, err := ah.Service.GetApplicationDeveloppers(appID)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   developpers,
	})
}

func (ah ApplicationHandler) GetApplicationConsumers(c *gin.Context) {
	appID := c.Param("applicationID")
	consumers, err := ah.Service.GetApplicationConsumers(appID)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   consumers,
	})
}

func (ah ApplicationHandler) List(c *gin.Context) {
	applications, err := ah.Service.GetAllApplications()
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   applications,
	})
}

func (ah ApplicationHandler) GetApplicationsCreatedByDevelopperID(c *gin.Context) {
	developperID := c.Param("developperID")

	applications, err := ah.Service.GetAllApplicationsCreatedByDevelopperID(developperID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   applications,
	})
}

func (ah ApplicationHandler) GetApplicationsCreatedForRoleID(c *gin.Context) {
	roleId := c.Param("roleID")

	applications, err := ah.Service.GetAllApplicationsCreatedForRoleID(roleId)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   applications,
	})
}

func (ah ApplicationHandler) AddApplicationConsumers(c *gin.Context) {

	userDevelopperID := getUserActor(c)
	var developperExists bool
	type roleConsumersDTO struct {
		ConsumersRoleID string `json:"consumer_id"`
	}
	var roleConsumerDtO roleConsumersDTO
	if err := c.ShouldBindJSON(&roleConsumerDtO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	applicationID := c.Param("applicationID")

	developpersForThisApplication, err := ah.Service.GetApplicationDeveloppers(applicationID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Erreur de traitement !",
		})
		return
	}
	developperExists = false
	for _, developper := range developpersForThisApplication {
		if developper.ID == userDevelopperID {
			developperExists = true
		}
	}

	if !developperExists {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Désole, vous n'êtes pas dévéloppeur de cette application",
		})
		return
	}

	applicationUpdated, err := ah.Service.AddConsumerRole(roleConsumerDtO.ConsumersRoleID, applicationID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// h.Logger.Info(fmt.Sprintf("Admin '%v' create a user with username = %v and id = %v", userActor.Username, user.Username, user.ID))
	//TODO: "à faire"
	c.JSON(201, gin.H{
		"status":    "success",
		"data":      applicationUpdated,
		"ressource": fmt.Sprintf("/applications/%v", applicationUpdated.ID),
	})
}

func (ah ApplicationHandler) AddApplicationDevelopper(c *gin.Context) {
	userDevelopperID := getUserActor(c)
	var developperExists bool
	type DevelopperDTO struct {
		DevelopperID string `json:"developper_id"`
	}
	var developperDTO DevelopperDTO
	if err := c.ShouldBindJSON(&developperDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	applicationID := c.Param("applicationID")

	developpersForThisApplication, err := ah.Service.GetApplicationDeveloppers(applicationID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Erreur de traitement !",
		})
		return
	}
	developperExists = false
	for _, developper := range developpersForThisApplication {
		if developper.ID == userDevelopperID {
			developperExists = true
		}
	}

	if !developperExists {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Désole, vous n'êtes pas dévéloppeur de cette application",
		})
		return
	}

	applicationUpdated, err := ah.Service.AddDevelopper(developperDTO.DevelopperID, applicationID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// h.Logger.Info(fmt.Sprintf("Admin '%v' create a user with username = %v and id = %v", userActor.Username, user.Username, user.ID))
	//TODO: "à faire"
	c.JSON(201, gin.H{
		"status":    "success",
		"data":      applicationUpdated,
		"ressource": fmt.Sprintf("/applications/%v", applicationUpdated.ID),
	})
}

func (ah ApplicationHandler) DeleteApplication(c *gin.Context) {

	userDevelopperID := getUserActor(c)
	var developperExists bool

	applicationID := c.Param("applicationID")

	developpersForThisApplication, err := ah.Service.GetApplicationDeveloppers(applicationID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Erreur de traitement !",
		})
		return
	}
	developperExists = false
	for _, developper := range developpersForThisApplication {
		if developper.ID == userDevelopperID {
			developperExists = true
		}
	}

	if !developperExists {
		c.JSON(401, gin.H{
			"status":  "error",
			"message": "Désole, vous n'êtes pas dévéloppeur de cette application",
		})
		return
	}

	err = ah.Service.DeleteApplication(applicationID)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "une erreur lors de la suppression de l'application",
			"details": err.Error(),
		})
		return
	}
	//TODO: "à faire"
	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Cette application a été supprimé avec succes",
	})
}
