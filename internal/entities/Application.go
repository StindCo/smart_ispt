package entities

type Application struct {
	ID            string
	Name          string
	ConsumerRoles []*Role
	Developpers   []*User
	PowerBy       string
	SmartName     string
	DomainName    string
	TestPath      string
	UrlPath       string
	Ip            string
	Description   string
}

func (a Application) AddConsumerRole(role *Role) {
	a.ConsumerRoles = append(a.ConsumerRoles, role)
}

func (a Application) RemoveConsumerRole(role *Role) {

}
