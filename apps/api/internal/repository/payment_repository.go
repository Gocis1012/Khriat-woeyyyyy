package repository

import (
	"context"
	"corporate-translator-api/internal/model"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrPaymentNotFound = errors.New("payment not found")

// pgxTxPool is the slice of *pgxpool.Pool this repository depends on.
// Both *pgxpool.Pool and pgxmock's pool satisfy it, so tests can inject a mock.
type pgxTxPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type PaymentRepository interface {
	CreatePending(ctx context.Context, userID string, amount float64, currency, method, providerTxID string) (*model.Payment, error)
	FindByID(ctx context.Context, id string) (*model.Payment, error)
	FindByProviderTxID(ctx context.Context, providerTxID string) (*model.Payment, error)
	// MarkProcessed marks a payment as processed and, on success, atomically
	// records a credit_ledger entry and credits the user's balance.
	MarkProcessed(ctx context.Context, paymentID, status string, creditDelta float64) error
}

type paymentRepository struct {
	db pgxTxPool
}

func NewPaymentRepository(db pgxTxPool) PaymentRepository {
	return &paymentRepository{db: db}
}

const paymentColumns = `id, user_id, amount, currency, method, status, provider_tx_id, webhook_processed_at, created_at`

func scanPayment(row pgx.Row) (*model.Payment, error) {
	var p model.Payment
	err := row.Scan(
		&p.ID, &p.UserID, &p.Amount, &p.Currency, &p.Method, &p.Status,
		&p.ProviderTxID, &p.WebhookProcessedAt, &p.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) CreatePending(ctx context.Context, userID string, amount float64, currency, method, providerTxID string) (*model.Payment, error) {
	query := fmt.Sprintf(`
	INSERT INTO payments (user_id, amount, currency, method, status, provider_tx_id)
	VALUES ($1, $2, $3, $4, 'pending', $5)
	RETURNING %s`, paymentColumns)

	row := r.db.QueryRow(ctx, query, userID, amount, currency, method, providerTxID)
	payment, err := scanPayment(row)
	if err != nil {
		return nil, fmt.Errorf("CreatePending: %w", err)
	}
	return payment, nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*model.Payment, error) {
	query := fmt.Sprintf(`SELECT %s FROM payments WHERE id = $1`, paymentColumns)

	payment, err := scanPayment(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return payment, nil
}

func (r *paymentRepository) FindByProviderTxID(ctx context.Context, providerTxID string) (*model.Payment, error) {
	query := fmt.Sprintf(`SELECT %s FROM payments WHERE provider_tx_id = $1`, paymentColumns)

	payment, err := scanPayment(r.db.QueryRow(ctx, query, providerTxID))
	if err != nil {
		if errors.Is(err, ErrPaymentNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("FindByProviderTxID: %w", err)
	}
	return payment, nil
}

func (r *paymentRepository) MarkProcessed(ctx context.Context, paymentID, status string, creditDelta float64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("MarkProcessed begin: %w", err)
	}
	defer tx.Rollback(ctx) // no-op once committed

	var userID string
	err = tx.QueryRow(ctx,
		`UPDATE payments SET status = $1, webhook_processed_at = now()
		 WHERE id = $2 RETURNING user_id`,
		status, paymentID,
	).Scan(&userID)
	if err != nil {
		return fmt.Errorf("MarkProcessed update payment: %w", err)
	}

	if status == "success" && creditDelta != 0 {
		_, err = tx.Exec(ctx,
			`INSERT INTO credit_ledger (user_id, amount, type, ref_id, ref_type)
			 VALUES ($1, $2, 'topup', $3, 'payment')`,
			userID, creditDelta, paymentID,
		)
		if err != nil {
			return fmt.Errorf("MarkProcessed insert ledger: %w", err)
		}

		_, err = tx.Exec(ctx,
			`UPDATE users SET credit = credit + $1, updated_at = now() WHERE id = $2`,
			creditDelta, userID,
		)
		if err != nil {
			return fmt.Errorf("MarkProcessed update user credit: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("MarkProcessed commit: %w", err)
	}
	return nil
}
