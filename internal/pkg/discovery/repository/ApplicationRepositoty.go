package repository

import (
	"github.com/StindCo/smart_ispt/internal/entities"
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
	SmartName     string `gorm:"unique;"`
	DomainName    string `gorm:"unique;"`
	TestPath      string
	UrlPath       string `gorm:"unique;"`
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
	app.UrlPath = entityApp.UrlPath
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
	var developpers []*identityRepo.UserGORM

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
	var consumerRoles []*identityRepo.RoleGORM

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

// TODO: Cette fonction a cruellement besoin d'un refactoring
func (a *ApplicationRepository) GetAllApplicationsByDeveloppers(developperID string) ([]*entities.Application, error) {
	type result struct {
		ApplicationID string
		UserGormID    string
	}

	developper := identityRepo.UserGORM{
		ID: developperID,
	}
	var results []*result

	a.DB.Model(&developper).Select("application_developpers.application_id, application_developpers.user_gorm_id").Joins("left join application_developpers on application_developpers.user_gorm_id = users.id").Scan(&results)
	applicationEntities := make([]*entities.Application, 0, len(results))

	for _, result := range results {
		app, _ := a.Get(result.ApplicationID)
		applicationEntities = append(applicationEntities, app)
	}
	return applicationEntities, nil
}

// TODO: Cette fonction a cruellement besoin d'un refactoring
func (a *ApplicationRepository) GetAllApplicationsByConsumers(consumerID string) ([]*entities.Application, error) {
	type result struct {
		ApplicationID string
		RoleGormID    string
	}

	consumer := identityRepo.RoleGORM{
		ID: consumerID,
	}
	var results []*result

	a.DB.Model(&consumer).Select("application_consumers.application_id, application_consumers.role_gorm_id").Joins("left join application_consumers on application_consumers.role_gorm_id = roles.id").Scan(&results)
	applicationEntities := make([]*entities.Application, 0, len(results))

	for _, result := range results {
		app, _ := a.Get(result.ApplicationID)
		applicationEntities = append(applicationEntities, app)
	}
	return applicationEntities, nil
}

func (a *ApplicationRepository) AddConsumerRole(application *entities.Application, role *entities.Role) (*entities.Application, error) {
	appGorm := NewApplicationGORM(application)
	roleGORM := identityRepo.NewRoleGORM(role)

	a.DB.Model(&appGorm).Association("ConsumerRoles").Append(roleGORM)

	return appGorm.toEntitiesApplication(), nil
}
