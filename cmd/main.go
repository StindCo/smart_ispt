package main

import (
	"github.com/StindCo/smart_ispt/pkg/database"
	"github.com/StindCo/smart_ispt/pkg/router"
)

func main() {

	db := database.RunConnectionToGorm()

	ginRouter := router.InitGinRouter()

	DispatchIdentityRouter(ginRouter.Group("/identity/v1"), db)
	DispatcheDiscoveryRouter(ginRouter.Group("/discovery"), db)

	router.Run(ginRouter)
}
