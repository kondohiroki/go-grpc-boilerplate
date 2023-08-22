package repository

import (
	"github.com/kondohiroki/go-grpc-boilerplate/internal/db/pgx"
)

type Repository struct {
	User UserRepository
	Job  JobRepository
}

func NewRepository() *Repository {
	readPgxPool := pgx.GetReadPgxPool()
	writePgxPool := pgx.GetWritePgxPool()

	return &Repository{
		User: NewUserRepository(readPgxPool, writePgxPool),
		Job:  NewJobRepository(readPgxPool, writePgxPool),
	}
}
