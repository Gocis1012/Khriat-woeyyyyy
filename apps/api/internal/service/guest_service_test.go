package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository"
	"errors"
	"testing"
)

// fakeGuestRepo implements repository.GuestRepository for tests.
type fakeGuestRepo struct {
	getOrCreate func(ctx context.Context, id string) (*model.Guest, error)
	deduct      func(ctx context.Context, id string, amt float64) error
	refund      func(ctx context.Context, id string, amt float64) error
	del         func(ctx context.Context, id string) error
}

func (f *fakeGuestRepo) GetOrCreate(ctx context.Context, id string) (*model.Guest, error) {
	return f.getOrCreate(ctx, id)
}
func (f *fakeGuestRepo) DeductCredit(ctx context.Context, id string, amt float64) error {
	return f.deduct(ctx, id, amt)
}
func (f *fakeGuestRepo) RefundCredit(ctx context.Context, id string, amt float64) error {
	return f.refund(ctx, id, amt)
}
func (f *fakeGuestRepo) Delete(ctx context.Context, id string) error {
	return f.del(ctx, id)
}

func TestGuestService_GetStatus(t *testing.T) {
	want := &model.Guest{Credit: 6}
	svc := NewGuestService(&fakeGuestRepo{
		getOrCreate: func(_ context.Context, _ string) (*model.Guest, error) { return want, nil },
	})
	got, err := svc.GetStatus(context.Background(), "g1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Credit != 6 {
		t.Errorf("credit = %v", got.Credit)
	}
}

func TestGuestService_GetStatus_Error(t *testing.T) {
	svc := NewGuestService(&fakeGuestRepo{
		getOrCreate: func(_ context.Context, _ string) (*model.Guest, error) {
			return nil, errors.New("boom")
		},
	})
	if _, err := svc.GetStatus(context.Background(), "g1"); err == nil {
		t.Error("expected error")
	}
}

func TestGuestService_UseCredit_Success(t *testing.T) {
	called := false
	svc := NewGuestService(&fakeGuestRepo{
		deduct: func(_ context.Context, _ string, _ float64) error { return nil },
	})
	err := svc.UseCredit(context.Background(), "g1", func() error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("callback was not invoked")
	}
}

func TestGuestService_UseCredit_InsufficientPassthrough(t *testing.T) {
	svc := NewGuestService(&fakeGuestRepo{
		deduct: func(_ context.Context, _ string, _ float64) error {
			return repository.ErrInsufficientCredit
		},
	})
	err := svc.UseCredit(context.Background(), "g1", func() error { return nil })
	if !errors.Is(err, repository.ErrInsufficientCredit) {
		t.Errorf("want ErrInsufficientCredit, got %v", err)
	}
}

func TestGuestService_UseCredit_GuestNotFoundPassthrough(t *testing.T) {
	svc := NewGuestService(&fakeGuestRepo{
		deduct: func(_ context.Context, _ string, _ float64) error {
			return repository.ErrGuestNotFound
		},
	})
	err := svc.UseCredit(context.Background(), "g1", func() error { return nil })
	if !errors.Is(err, repository.ErrGuestNotFound) {
		t.Errorf("want ErrGuestNotFound, got %v", err)
	}
}

func TestGuestService_UseCredit_DeductOtherError(t *testing.T) {
	svc := NewGuestService(&fakeGuestRepo{
		deduct: func(_ context.Context, _ string, _ float64) error { return errors.New("redis down") },
	})
	err := svc.UseCredit(context.Background(), "g1", func() error { return nil })
	if err == nil {
		t.Error("expected error")
	}
}

func TestGuestService_UseCredit_RefundsOnCallbackFailure(t *testing.T) {
	refunded := false
	svc := NewGuestService(&fakeGuestRepo{
		deduct: func(_ context.Context, _ string, _ float64) error { return nil },
		refund: func(_ context.Context, _ string, _ float64) error { refunded = true; return nil },
	})
	err := svc.UseCredit(context.Background(), "g1", func() error {
		return errors.New("AI failed")
	})
	if err == nil {
		t.Error("expected error from callback")
	}
	if !refunded {
		t.Error("expected refund to be called")
	}
}

func TestGuestService_UseCredit_RefundFailureStillReturnsCallbackError(t *testing.T) {
	svc := NewGuestService(&fakeGuestRepo{
		deduct: func(_ context.Context, _ string, _ float64) error { return nil },
		refund: func(_ context.Context, _ string, _ float64) error { return errors.New("refund failed") },
	})
	err := svc.UseCredit(context.Background(), "g1", func() error {
		return errors.New("AI failed")
	})
	if err == nil {
		t.Error("expected callback error even when refund fails")
	}
}

func TestGuestService_DeleteSession(t *testing.T) {
	called := false
	svc := NewGuestService(&fakeGuestRepo{
		del: func(_ context.Context, _ string) error { called = true; return nil },
	})
	if err := svc.DeleteSession(context.Background(), "g1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected Delete to be called")
	}
}
