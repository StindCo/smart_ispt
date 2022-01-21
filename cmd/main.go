package main

import (
	"github.com/StindCo/smart_ispt/pkg/database"
	"github.com/StindCo/smart_ispt/pkg/router"
)

func main() {
	db := database.RunConnectionToGorm()

	ginRouter := router.InitGinRouter()

	LaunchIdentity(ginRouter.Group("/identity"), db)

	router.Run(ginRouter)
}
