package db

import (
	"context"

	"github.com/lookeme/short-url/internal/models"
)

type UserRepository struct {
	postgres *Postgres
}

func NewUserRepository(postgres *Postgres) *UserRepository {
	return &UserRepository{
		postgres: postgres,
	}
}

func (u *UserRepository) SaveUser(name, pass string) (int, error) {
	lastInsertID := 0
	err := u.postgres.connPool.QueryRow(
		context.Background(),
		"INSERT INTO users(name, pass) VALUES($1, $2) RETURNING id",
		name, pass).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, err
	}
	return lastInsertID, nil
}
func (u *UserRepository) FindByID(userID int) (models.User, error) {
	return models.User{}, nil
}
