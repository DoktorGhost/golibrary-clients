package providers

import (
	"github.com/DoktorGhost/golibrary-clients/internal/usecases"
)

type UseCaseProvider struct {
	UserUseCase *usecases.UsersUseCase
}

func NewUseCaseProvider() *UseCaseProvider {
	return &UseCaseProvider{}
}

func (ucp *UseCaseProvider) RegisterDependencies(provider *ServiceProvider) {
	ucp.UserUseCase = usecases.NewUsersUseCase(provider.usersService)

}
