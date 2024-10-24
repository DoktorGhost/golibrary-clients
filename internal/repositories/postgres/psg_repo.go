package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *UsersRepository {
	return &UsersRepository{db: db}
}
