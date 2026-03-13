package repository

import (
	"Web/user-service/internal/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) (int, error)
	GetByFirstName(ctx context.Context, login string) (*models.User, error)
	GetById(ctx context.Context, id int) (*models.User, error)
}
