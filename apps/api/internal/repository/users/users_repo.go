package users

import (
	"context"
	"corporate-translator-api/internal/model"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Insert(ctx context.Context, user *model.User) error
	FindByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	DeductCredit(ctx context.Context, id string, amount float64) error
}

type postgrestRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) UserRepository {
	return &postgrestRepository{db: db}
}

func (r *postgrestRepository) Insert(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users (google_id, email, username, avatar_url)
	VALUES ($1, $2, $3, $4)
	RETURNING id, credit, member_type, created_at, updated_at`

	row := r.db.QueryRow(ctx, query,
		user.GoogleID,
		user.Email,
		user.Username,
		user.AvatarURL,
	)

	err := row.Scan(
		&user.ID,
		&user.Credit,
		&user.MemberType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert repo: %w", err)
	}
	return nil
}

func (r *postgrestRepository) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	query := `
	SELECT id, google_id, email, username, avatar_url, credit, member_type,
	       last_daily_credit_at, created_at, updated_at
	FROM users WHERE google_id = $1`

	var user model.User
	err := r.db.QueryRow(ctx, query, googleID).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Username, &user.AvatarURL,
		&user.Credit, &user.MemberType, &user.LastDailyCreditAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("FindByGoogleID: %w", err)
	}
	return &user, nil
}

func (r *postgrestRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
	SELECT id, google_id, email, username, avatar_url, credit, member_type,
	       last_daily_credit_at, created_at, updated_at
	FROM users WHERE id = $1`

	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Username, &user.AvatarURL,
		&user.Credit, &user.MemberType, &user.LastDailyCreditAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return &user, nil
}

func (r *postgrestRepository) DeductCredit(ctx context.Context, id string, amount float64) error {
	query := `
	UPDATE users SET credit = credit - $1, updated_at = now()
	WHERE id = $2 AND credit >= $1
	RETURNING id`

	var returnedID string
	err := r.db.QueryRow(ctx, query, amount, id).Scan(&returnedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("DeductCredit: insufficient credit or user not found")
	}
	if err != nil {
		return fmt.Errorf("DeductCredit: %w", err)
	}
	return nil
}
