package router

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func InitGinRouter() *gin.Engine {

	gin.ForceConsoleColor()
	app := gin.Default()
	app.SetTrustedProxies([]string{"192.168.1.2"})
	app.Use(gin.Recovery())

	return app
}

func Run(app *gin.Engine) {

	port := os.Getenv("PORT")
	if port == "" {
		port = "4500"
	}

	err := app.Run(fmt.Sprintf(":%v", port))
	if err != nil {
		fmt.Println("Il y'a erreur lors du démarrage du serveur")
	}
}
