package service

import (
	"errors"

	"github.com/StindCo/smart_ispt/internal/entities"
	"github.com/StindCo/smart_ispt/internal/pkg/discovery/interfaces"
	discoveryRepo "github.com/StindCo/smart_ispt/internal/pkg/discovery/repository"
	repository "github.com/StindCo/smart_ispt/internal/pkg/identity/repository"
)

type ApplicationServiceImpl struct {
	UserRepository        repository.UserRepository
	RoleRepository        repository.RoleRepository
	ApplicationRepository discoveryRepo.ApplicationRepository
}

func NewApplicationService(ar discoveryRepo.ApplicationRepository, ur repository.UserRepository, rr repository.RoleRepository) interfaces.ApplicationService {
	return &ApplicationServiceImpl{
		UserRepository:        ur,
		RoleRepository:        rr,
		ApplicationRepository: ar,
	}
}

func (as ApplicationServiceImpl) CreateApplication(developperID string, entityApp *entities.Application) (*entities.Application, error) {
	user, err := as.UserRepository.Get(developperID)
	if err != nil {
		return nil, err
	}

	if user.IsDevelopper == 0 {
		return nil, errors.New("désolé, vous n'êtes pas administrateur pour effectuer cette action")
	}

	entityApp.Developpers = append(entityApp.Developpers, user)
	applicationCreated, err := as.ApplicationRepository.Create(entityApp)
	if err != nil {
		return nil, err
	}

	return applicationCreated, nil
}

func (as ApplicationServiceImpl) GetOneApplication(applicationID string) (*entities.Application, error) {
	application, err := as.ApplicationRepository.Get(applicationID)
	if err != nil {
		return nil, errors.New("désolé, une erreur interne")
	}
	return application, err
}

func (as ApplicationServiceImpl) GetApplicationDeveloppers(applicationID string) ([]*entities.User, error) {
	developpers, err := as.ApplicationRepository.GetDeveloppersForApplicationID(applicationID)
	if err != nil {
		return nil, errors.New("désolé, une erreur interne")
	}
	return developpers, nil
}

func (as ApplicationServiceImpl) GetApplicationConsumers(applicationID string) ([]*entities.Role, error) {
	consumers, err := as.ApplicationRepository.GetConsumersForApplicationID(applicationID)
	if err != nil {
		return nil, errors.New("désolé, une erreur interne")
	}
	return consumers, nil
}

func (as ApplicationServiceImpl) GetAllApplications() ([]*entities.Application, error) {
	applications, err := as.ApplicationRepository.GetAllApplications()
	if err != nil {
		return nil, errors.New("désolé, une erreur interne")
	}
	return applications, nil
}
