package service

import (
	"Web/user-service/internal/auth"
	"Web/user-service/internal/models"
	"Web/user-service/internal/repository"
	"context"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
	jwt  *auth.JWTManager
}

func NewUserService(r repository.UserRepository, jwt *auth.JWTManager) *UserService {
	return &UserService{
		repo: r,
		jwt:  jwt,
	}
}

func (s *UserService) Register(ctx context.Context, first, last, birth, password, login string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	birthDate, err := time.Parse("2006-01-02", birth)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Login:        login,
		FirstName:    first,
		LastName:     last,
		BirthDate:    birthDate,
		PasswordHash: string(hash),
	}

	id, err := s.repo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return s.jwt.Generate(id)
}

func (s *UserService) Login(ctx context.Context, login, password string) (string, error) {

	user, err := s.repo.GetByFirstName(ctx, login)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	log.Println("LOGIN NAME:", login)

	return s.jwt.Generate(user.ID)
}

func (s *UserService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	return s.repo.GetById(ctx, userID)
}
