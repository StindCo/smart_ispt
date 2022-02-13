package handler

import (
	"strings"

	"github.com/StindCo/smart_ispt/internal/pkg/discovery/interfaces"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-hclog"
)

type RunnerHandler struct {
	Service interfaces.ApplicationService
	Logger  hclog.Logger
}

func NewRunnerHandler(app *gin.RouterGroup, auth *jwt.GinJWTMiddleware, service interfaces.ApplicationService, logger hclog.Logger) {
	runnerHandler := RunnerHandler{
		Service: service,
		Logger:  logger,
	}
	app.Any(":smartName", runnerHandler.Run)
	app.Any(":smartName/*action", runnerHandler.Run)
}

func (rh RunnerHandler) Run(c *gin.Context) {
	smartName := c.Param("smartName")
	action := c.Param("action")

	application, err := rh.Service.GetOneApplicationBySmartName(smartName)

	if err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "erreur interne",
			"details": err.Error(),
		})
	}

	client := resty.New()
	request := client.R()

	data, _ := c.GetRawData()
	request.Header = c.Request.Header
	request.SetBody(data)

	url := application.UrlPath + action
	resp, _ := request.Execute(strings.ToUpper(c.Request.Method), url)
	c.Request.Response = resp.RawResponse
	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
}
