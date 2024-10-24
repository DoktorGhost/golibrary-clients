package providers

import (
	"github.com/DoktorGhost/golibrary-clients/internal/services"
)

type ServiceProvider struct {
	usersService *services.UserService
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) RegisterDependencies(provider *RepositoryProvider) {
	s.usersService = services.NewUserService(provider.usersRepositoryPostgres)
}
