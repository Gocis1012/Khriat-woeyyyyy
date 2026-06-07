package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"errors"
	"testing"
)

// fakeUserRepo implements users.UserRepository for tests.
type fakeUserRepo struct {
	insert         func(ctx context.Context, u *model.User) error
	findByGoogleID func(ctx context.Context, gid string) (*model.User, error)
	findByID       func(ctx context.Context, id string) (*model.User, error)
	deduct         func(ctx context.Context, id string, amt float64) error
}

func (f *fakeUserRepo) Insert(ctx context.Context, u *model.User) error { return f.insert(ctx, u) }
func (f *fakeUserRepo) FindByGoogleID(ctx context.Context, gid string) (*model.User, error) {
	return f.findByGoogleID(ctx, gid)
}
func (f *fakeUserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	return f.findByID(ctx, id)
}
func (f *fakeUserRepo) DeductCredit(ctx context.Context, id string, amt float64) error {
	return f.deduct(ctx, id, amt)
}

func TestUserService_Insert_Validation(t *testing.T) {
	svc := NewUserService(&fakeUserRepo{
		insert: func(_ context.Context, _ *model.User) error { return nil },
	})

	if err := svc.Insert(context.Background(), nil); err == nil {
		t.Error("expected error for nil user")
	}
	if err := svc.Insert(context.Background(), &model.User{Email: ""}); err == nil {
		t.Error("expected error for empty email")
	}
}

func TestUserService_Insert_Success(t *testing.T) {
	called := false
	svc := NewUserService(&fakeUserRepo{
		insert: func(_ context.Context, _ *model.User) error { called = true; return nil },
	})
	err := svc.Insert(context.Background(), &model.User{Email: "a@b.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("repo.Insert not called")
	}
}

func TestUserService_Insert_RepoError(t *testing.T) {
	svc := NewUserService(&fakeUserRepo{
		insert: func(_ context.Context, _ *model.User) error { return errors.New("db error") },
	})
	if err := svc.Insert(context.Background(), &model.User{Email: "a@b.com"}); err == nil {
		t.Error("expected wrapped error")
	}
}

func TestUserService_FindByGoogleID(t *testing.T) {
	want := &model.User{ID: "u1"}
	svc := NewUserService(&fakeUserRepo{
		findByGoogleID: func(_ context.Context, _ string) (*model.User, error) { return want, nil },
	})
	got, err := svc.FindByGoogleID(context.Background(), "gid")
	if err != nil || got.ID != "u1" {
		t.Errorf("got %v, err %v", got, err)
	}
}

func TestUserService_FindByID(t *testing.T) {
	want := &model.User{ID: "u2"}
	svc := NewUserService(&fakeUserRepo{
		findByID: func(_ context.Context, _ string) (*model.User, error) { return want, nil },
	})
	got, err := svc.FindByID(context.Background(), "u2")
	if err != nil || got.ID != "u2" {
		t.Errorf("got %v, err %v", got, err)
	}
}

func TestUserService_DeductCredit(t *testing.T) {
	called := false
	svc := NewUserService(&fakeUserRepo{
		deduct: func(_ context.Context, _ string, amt float64) error {
			called = true
			if amt != 1.0 {
				t.Errorf("amount = %v", amt)
			}
			return nil
		},
	})
	if err := svc.DeductCredit(context.Background(), "u1", 1.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("repo.DeductCredit not called")
	}
}
