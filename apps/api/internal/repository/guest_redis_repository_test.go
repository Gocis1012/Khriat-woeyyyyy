package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func newTestRepo(t *testing.T) (GuestRepository, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return NewGuestRepository(client), mr
}

func TestGuestRepo_GetOrCreate_New(t *testing.T) {
	repo, _ := newTestRepo(t)
	g, err := repo.GetOrCreate(context.Background(), "g1")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}
	if g.Credit != guestInitialCredit {
		t.Errorf("credit = %v, want %v", g.Credit, guestInitialCredit)
	}
}

func TestGuestRepo_GetOrCreate_Existing(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	_, _ = repo.GetOrCreate(ctx, "g1")
	_ = repo.DeductCredit(ctx, "g1", 2.0)

	g, err := repo.GetOrCreate(ctx, "g1")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}
	if g.Credit != guestInitialCredit-2.0 {
		t.Errorf("credit = %v, want %v", g.Credit, guestInitialCredit-2.0)
	}
}

func TestGuestRepo_DeductCredit(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	_, _ = repo.GetOrCreate(ctx, "g1")

	if err := repo.DeductCredit(ctx, "g1", 1.0); err != nil {
		t.Fatalf("DeductCredit: %v", err)
	}
	g, _ := repo.GetOrCreate(ctx, "g1")
	if g.Credit != guestInitialCredit-1.0 {
		t.Errorf("credit = %v", g.Credit)
	}
}

func TestGuestRepo_DeductCredit_Insufficient(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	_, _ = repo.GetOrCreate(ctx, "g1")

	err := repo.DeductCredit(ctx, "g1", 9999.0)
	if !errors.Is(err, ErrInsufficientCredit) {
		t.Errorf("want ErrInsufficientCredit, got %v", err)
	}
}

func TestGuestRepo_DeductCredit_GuestNotFound(t *testing.T) {
	repo, _ := newTestRepo(t)
	err := repo.DeductCredit(context.Background(), "ghost", 1.0)
	if !errors.Is(err, ErrGuestNotFound) {
		t.Errorf("want ErrGuestNotFound, got %v", err)
	}
}

func TestGuestRepo_RefundCredit(t *testing.T) {
	repo, _ := newTestRepo(t)
	ctx := context.Background()
	_, _ = repo.GetOrCreate(ctx, "g1")
	_ = repo.DeductCredit(ctx, "g1", 3.0)

	if err := repo.RefundCredit(ctx, "g1", 2.0); err != nil {
		t.Fatalf("RefundCredit: %v", err)
	}
	g, _ := repo.GetOrCreate(ctx, "g1")
	if g.Credit != guestInitialCredit-3.0+2.0 {
		t.Errorf("credit = %v", g.Credit)
	}
}

func TestGuestRepo_RefundCredit_NoKeyIsNoop(t *testing.T) {
	repo, _ := newTestRepo(t)
	if err := repo.RefundCredit(context.Background(), "ghost", 1.0); err != nil {
		t.Errorf("refund on missing key should be a no-op, got %v", err)
	}
}

func TestGuestRepo_Delete(t *testing.T) {
	repo, mr := newTestRepo(t)
	ctx := context.Background()
	_, _ = repo.GetOrCreate(ctx, "g1")

	if err := repo.Delete(ctx, "g1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if mr.Exists("guest:g1") {
		t.Error("key should be deleted")
	}
}
