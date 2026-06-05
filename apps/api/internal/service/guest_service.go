package service

import (
	"context"
	"corporate-translator-api/internal/model"
	"corporate-translator-api/internal/repository"
	"errors"
	"fmt"
	"log/slog"
)

const guestCreditCost = 1.00

type GuestService struct {
	repo repository.GuestRepository
}

func NewGuestService(repo repository.GuestRepository) *GuestService {
	return &GuestService{repo: repo}
}

func (s *GuestService) GetStatus(ctx context.Context, guestID string) (*model.Guest, error) {

	guest, err := s.repo.GetOrCreate(ctx, guestID)

	if err != nil {
		return nil, fmt.Errorf("GuestService.GetStatus: %w", err)

	}

	return guest, nil
}

func (s *GuestService) UseCredit(ctx context.Context, guestID string, fn func() error) error {

	err := s.repo.DeductCredit(ctx, guestID, guestCreditCost)

	if err != nil {
		// errors.Is() เช็คว่า error ตัวนี้คือ ErrInsufficientCredit ไหม
		// แม้จะถูก wrap ด้วย fmt.Errorf("%w") มาหลายชั้น ก็ยังเจอ
		if errors.Is(err, repository.ErrInsufficientCredit) {
			return err // ส่งต่อให้ handler จัดการ
		}
		if errors.Is(err, repository.ErrGuestNotFound) {
			return err
		}
		return fmt.Errorf("GuestService.UseCredit deduct: %w", err)

	}

	if err := fn(); err != nil {
		refundErr := s.repo.RefundCredit(ctx, guestID, guestCreditCost)
		if refundErr != nil {
			// log ไว้ แต่ return error ของ fn เป็นหลัก
			slog.Warn("Refund failed for guest", "guest_id", guestID, "error", refundErr)
		}
		return fmt.Errorf("GuestService.UseCredit fn: %w", err)
	}

	return nil
}


func (s *GuestService) DeleteSession(ctx context.Context, guestID string) error {
	return s.repo.Delete(ctx, guestID)
}

