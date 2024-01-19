package repository

import (
	"testApp/internal/repository/users"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	Users users.UserRepository
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{
		Users: users.NewUserRepository(db),
	}
}
