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

func (as ApplicationServiceImpl) GetOneApplicationBySmartName(smartName string) (*entities.Application, error) {
	application, err := as.ApplicationRepository.GetBySmartName(smartName)
	if err != nil {
		return nil, err
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

func (as ApplicationServiceImpl) GetAllApplicationsCreatedByDevelopperID(developperID string) ([]*entities.Application, error) {
	user, err := as.UserRepository.Get(developperID)
	if err != nil {
		return nil, errors.New("désolé, cet utlisateur n'existe pas")
	}
	if user.IsDevelopper == 0 {
		return nil, errors.New("cet utilisateur, n'est pas présentement dévéloppeur, donc n'a pas d'application à son actif")
	}
	applications, err := as.ApplicationRepository.GetAllApplicationsByDeveloppers(developperID)

	return applications, err
}

func (as ApplicationServiceImpl) GetAllApplicationsCreatedForRoleID(roleID string) ([]*entities.Application, error) {
	_, err := as.UserRepository.Get(roleID)
	if err != nil {
		return nil, errors.New("désolé, ce role n'existe pas")
	}
	applications, err := as.ApplicationRepository.GetAllApplicationsByDeveloppers(roleID)

	return applications, err
}

func (as ApplicationServiceImpl) AddConsumerRole(roleID string, applicationID string) (*entities.Application, error) {
	role, err := as.RoleRepository.Get(roleID)
	if err != nil {
		return nil, errors.New("ce role n'existe pas")
	}
	application, err := as.ApplicationRepository.Get(applicationID)
	if err != nil {
		return nil, errors.New("cette application n'existe pas, sorry ")
	}
	application, err = as.ApplicationRepository.AddConsumerRole(application, role)
	return application, err
}

func (as ApplicationServiceImpl) AddDevelopper(developperID string, applicationID string) (*entities.Application, error) {
	developper, err := as.UserRepository.Get(developperID)
	if err != nil {
		return nil, errors.New("cet utilisateur n'existe pas ")
	}
	if developper.IsDevelopper == 0 {
		return nil, errors.New("cet utilisateur n'est pas dévéloppeur")
	}
	application, err := as.ApplicationRepository.Get(applicationID)
	if err != nil {
		return nil, errors.New("cette application n'existe pas, sorry ")
	}
	application, err = as.ApplicationRepository.AddDevelopper(application, developper)
	return application, err
}

func (as ApplicationServiceImpl) DeleteApplication(applicationID string) error {
	return as.ApplicationRepository.DeleteApplication(applicationID)
}
