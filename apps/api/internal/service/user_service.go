package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository/users"
	"fmt"
)

type UserService interface {
	Insert(ctx context.Context, user *model.User) error
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	DeductCredit(ctx context.Context, id string, amount float64) error
}

type userService struct {
	repo users.UserRepository
}

func NewUserService(repo users.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Insert(ctx context.Context, user *model.User) error {
	if user == nil || user.Email == "" {
		return fmt.Errorf("invalid input")
	}
	if err := s.repo.Insert(ctx, user); err != nil {
		return fmt.Errorf("user service insert: %w", err)
	}
	return nil
}

func (s *userService) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	return s.repo.FindByGoogleID(ctx, googleID)
}

func (s *userService) FindByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) DeductCredit(ctx context.Context, id string, amount float64) error {
	return s.repo.DeductCredit(ctx, id, amount)
}
