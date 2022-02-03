package database

import (
	discoveryRepo "github.com/StindCo/smart_ispt/internal/pkg/discovery/repository"
	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
	"github.com/StindCo/smart_ispt/pkg/applogger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectGORMDB(dialector gorm.Dialector) (*gorm.DB, error) {
	var gormDB *gorm.DB
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		// Logger:                                   logger.Default.LogMode(logger.Silent),
		// DryRun: true,
	})
	if err != nil {
		return nil, err
	}

	gormDB.AutoMigrate(&repository.UserGORM{}, &repository.RoleGORM{}, &discoveryRepo.Application{})
	return gormDB, nil
}

func RunConnectionToGorm() *gorm.DB {
	dbLogger := applogger.NewLogger("database")

	dsn := "stephane:djodjo789+456@tcp(127.0.0.1:3306)/smart_ispt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := ConnectGORMDB(mysql.Open(dsn))
	if err != nil {
		dbLogger.Info("Une erreur lors de la connexion à la base de donnée")
	}
	return db
}
