package repository

import (
	"github.com/StindCo/smart_ispt/internal/entities"
	"github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
	identityRepo "github.com/StindCo/smart_ispt/internal/pkg/identity/repository"

	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	ID            string `gorm:"primarykey"`
	Name          string
	ConsumerRoles []*identityRepo.RoleGORM `gorm:"many2many:application_consumers;"`
	Developpers   []*identityRepo.UserGORM `gorm:"many2many:application_developpers;"`
	PowerBy       string
	SmartName     string
	DomainName    string
	TestPath      string
	UrlPath       string
	Ip            string
	Description   string
}

func (a Application) toEntitiesApplication() *entities.Application {
	developpers := make([]*entities.User, 0, len(a.Developpers))
	for _, developper := range a.Developpers {
		developpers = append(developpers, developper.ToEntitiesUser())
	}
	roles := make([]*entities.Role, 0, len(a.ConsumerRoles))
	for _, role := range a.ConsumerRoles {
		roles = append(roles, role.ToEntitiesRole())
	}
	return &entities.Application{
		ID:            a.ID,
		Developpers:   developpers,
		ConsumerRoles: roles,
		Name:          a.Name,
		UrlPath:       a.UrlPath,
		SmartName:     a.SmartName,
		PowerBy:       a.PowerBy,
		TestPath:      a.TestPath,
		Description:   a.Description,
		DomainName:    a.DomainName,
		Ip:            a.Ip,
	}
}

func NewApplicationGORM(entityApp *entities.Application) *Application {
	developpers := make([]*identityRepo.UserGORM, 0, len(entityApp.Developpers))
	for _, developper := range entityApp.Developpers {
		developpers = append(developpers, identityRepo.NewUserGORM(developper))
	}
	roles := make([]*identityRepo.RoleGORM, 0, len(entityApp.ConsumerRoles))
	for _, role := range entityApp.ConsumerRoles {
		roles = append(roles, identityRepo.NewRoleGORM(role))
	}

	app := Application{}
	app.ID = entityApp.ID
	app.DomainName = entityApp.DomainName
	app.Ip = entityApp.Ip
	app.TestPath = entityApp.TestPath
	app.Description = entityApp.Description
	app.PowerBy = entityApp.PowerBy
	app.SmartName = entityApp.SmartName
	app.Name = entityApp.Name
	app.Developpers = developpers
	app.ConsumerRoles = roles
	return &app
}

type ApplicationRepository struct {
	DB *gorm.DB
}

func NewApplicationGORMRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{
		DB: db,
	}
}

func (a *ApplicationRepository) Create(entityApp *entities.Application) (*entities.Application, error) {
	application := NewApplicationGORM(entityApp)
	err := a.DB.Omit("ConsumerRoles").Create(&application).Error
	return application.toEntitiesApplication(), err
}

func (a *ApplicationRepository) Get(applicationID string) (*entities.Application, error) {
	var application Application
	err := a.DB.Where("id = ?", applicationID).First(&application).Error
	if err != nil {
		return nil, err
	}
	return application.toEntitiesApplication(), err
}

func (a *ApplicationRepository) GetDeveloppersForApplicationID(applicationID string) ([]*entities.User, error) {
	var application Application
	application.ID = applicationID
	var developpers []*repository.UserGORM

	a.DB.Model(&application).Association("Developpers").Find(&developpers)

	developpersFinals := make([]*entities.User, 0, len(developpers))
	for _, developper := range developpers {
		developpersFinals = append(developpersFinals, developper.ToEntitiesUser())
	}

	return developpersFinals, nil
}

func (a *ApplicationRepository) GetConsumersForApplicationID(applicationID string) ([]*entities.Role, error) {
	var application Application
	application.ID = applicationID
	var consumerRoles []*repository.RoleGORM

	a.DB.Model(&application).Association("ConsumerRoles").Find(&consumerRoles)

	consumersRoles := make([]*entities.Role, 0, len(consumerRoles))
	for _, role := range consumerRoles {
		consumersRoles = append(consumersRoles, role.ToEntitiesRole())
	}

	return consumersRoles, nil
}

func (a *ApplicationRepository) GetAllApplications() ([]*entities.Application, error) {
	var applications []*Application

	a.DB.Model(&applications).Find(&applications)

	applicationEntities := make([]*entities.Application, 0, len(applications))
	for _, application := range applications {
		applicationEntities = append(applicationEntities, application.toEntitiesApplication())
	}

	return applicationEntities, nil
}
