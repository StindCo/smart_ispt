package database

import (
	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/Repository"
	"gorm.io/gorm"
)

func ConnectGORMDB(dialector gorm.Dialector) (*gorm.DB, error) {
	var gormDB *gorm.DB
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	gormDB.AutoMigrate(&repository.UserGORM{}, &repository.RoleGORM{})
	return gormDB, nil
}
