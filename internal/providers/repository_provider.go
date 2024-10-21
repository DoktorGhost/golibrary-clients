package providers

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DoktorGhost/golibrary-clients/internal/repositories/postgres"
)

type RepositoryProvider struct {
	db                      *pgxpool.Pool
	usersRepositoryPostgres *postgres.UsersRepository
}

func NewRepositoryProvider(db *pgxpool.Pool) *RepositoryProvider {
	return &RepositoryProvider{db: db}
}

func (r *RepositoryProvider) RegisterDependencies() {
	r.usersRepositoryPostgres = postgres.NewPostgresRepository(r.db)
}
