package app

import (
	"github.com/DoktorGhost/golibrary-clients/internal/providers"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

var (
	Container container
	once      sync.Once
)

type container struct {
	UseCaseProvider *providers.UseCaseProvider
}

func Init(db *pgxpool.Pool) container {
	once.Do(func() {
		repositoryProvider := providers.NewRepositoryProvider(db)
		repositoryProvider.RegisterDependencies()

		serviceProvider := providers.NewServiceProvider()
		serviceProvider.RegisterDependencies(repositoryProvider)

		useCaseProvider := providers.NewUseCaseProvider()
		useCaseProvider.RegisterDependencies(serviceProvider)

		Container = container{
			UseCaseProvider: useCaseProvider,
		}
	})
	return Container
}
