package app

import (
	"context"

	"github.com/kondohiroki/go-grpc-boilerplate/internal/repository"
)

// App is a business logic layer.
// It is used to separate the business logic layer from the interface layer.
type App interface {
	GetUsers(ctx context.Context) ([]GetUserDTO, error)
	GetUserByID(ctx context.Context, input GetUserDTI) (GetUserDTO, error)
	CreateUser(ctx context.Context, input CreateUserDTI) (CreateUserDTO, error)
}

type app struct {
	Repo *repository.Repository
}

func NewApp(repo *repository.Repository) App {
	return &app{
		Repo: repo,
	}
}
