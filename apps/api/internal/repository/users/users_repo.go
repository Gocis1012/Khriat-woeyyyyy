package users

import (
	"context"
	"corporate-translator-api/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Insert(ctx context.Context, users *model.User) error
}

type postgrestRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) UserRepository {
	return &postgrestRepository{db: db}
}

func (r *postgrestRepository) Insert(ctx context.Context, users *model.User) error {
    query := `
    INSERT INTO users (
        google_id,
        email,
        username,
        avatar_url
    ) VALUES (
        $1, $2, $3, $4
    ) RETURNING id, credit, member_type, created_at, updated_at`

    row := r.db.QueryRow(ctx, query,
        users.GoogleID,
        users.Email,
        users.Username,
        users.AvatarURL,
    )

    // รับค่าที่ DB generate กลับมาด้วย
    err := row.Scan(
        &users.ID,
        &users.Credit,
        &users.MemberType,
        &users.CreatedAt,
        &users.UpdatedAt,
    )
    if err != nil {
        return fmt.Errorf("insert repo: %w", err)
    }

    return nil
}