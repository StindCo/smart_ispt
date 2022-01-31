package router

import (
	"fmt"
	"io"
	"os"

	"github.com/StindCo/smart_ispt/pkg/applogger"
	"github.com/gin-gonic/gin"
)

func InitGinRouter() *gin.Engine {

	gin.DisableConsoleColor()
	// Logging to a file.
	f, err := os.Create("./log/gin.log")
	if err != nil {
		fmt.Println(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)

	app := gin.New()
	app.SetTrustedProxies([]string{"192.168.1.2"})
	app.Use(gin.Recovery())
	app.Use(gin.Logger())
	return app
}

func Run(app *gin.Engine) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "4500"
	}

	err := app.Run(fmt.Sprintf(":%v", port))
	if err == nil {
		applogger.NewLogger("Application Root").Info("Server started")
	}
	fmt.Println("Il y'a erreur lors du démarrage du serveur")
}
