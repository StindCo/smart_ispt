package main

import (
	"fmt"

	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/Repository"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/handler"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/service"
	"github.com/StindCo/smart_ispt/pkg/database"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
)

func main() {

	dsn := "stephane:djodjo789+456@tcp(127.0.0.1:3306)/smart_ispt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := database.ConnectGORMDB(mysql.Open(dsn))
	if err != nil {
		fmt.Println("Hello world")
	}

	app := gin.New()
	app.SetTrustedProxies([]string{"192.168.1.2"})
	gin.Recovery()

	userRepository := repository.NewUserGORMRepository(db)

	userService := service.NewUserService(*userRepository)

	api := app.Group("/api")

	handler.NewUserHandler(api.Group("v1/users"), userService)

	app.Run(":4500")
}
