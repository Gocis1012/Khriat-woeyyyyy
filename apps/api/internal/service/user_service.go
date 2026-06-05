package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository/users"
	"fmt"
)

type UserService interface {
	Insert(ctx context.Context, users *model.User) error
}

type userService struct {
	repo users.UserRepository
}

func NewUserService(repo users.UserRepository) UserService {
	return &userService{repo: repo}
}

func (r *userService) Insert(ctx context.Context, user *model.User) error {
	if user == nil || user.Email == "" {
		return fmt.Errorf("invalid Input")
	}

	if err := r.repo.Insert(ctx, user); err != nil {
        return fmt.Errorf("user service insert: %w", err)     // ✅ wrap error
    }

	return nil
}
