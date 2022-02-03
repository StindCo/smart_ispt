package interfaces

import "github.com/StindCo/smart_ispt/internal/entities"

type ApplicationService interface {
	CreateApplication(developperID string, entityApp *entities.Application) (*entities.Application, error)

	GetOneApplication(applicationID string) (*entities.Application, error)
	GetApplicationDeveloppers(applicationID string) ([]*entities.User, error)
	GetApplicationConsumers(applicationID string) ([]*entities.Role, error)
	GetAllApplications() ([]*entities.Application, error)

	GetAllApplicationsCreatedByDevelopperID(developperID string) ([]*entities.Application, error)
	GetAllApplicationsCreatedForRoleID(roleID string) ([]*entities.Application, error)

	// UpdateApplication() (*entities.Application, error)

	AddConsumerRole(roleID string, applicationID string) (*entities.Application, error)
	// RemoveConsumerRole(roleID string, applicationID string) (*entities.Application, error)

	// AddDevelopper(developperID string, applicationID string) (*entities.Application, error)

	// DeleteApplication(applicationID string) error
}
