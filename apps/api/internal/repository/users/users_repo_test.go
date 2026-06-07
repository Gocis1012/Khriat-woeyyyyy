package users

import (
	"context"
	"errors"
	"testing"
	"time"

	"corporate-translator-api/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func TestUsersRepo_Insert(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "credit", "member_type", "created_at", "updated_at"}).
		AddRow("user-1", 10.0, "free", now, now)
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("gid", "e@x.com", "name", (*string)(nil), 10.0).
		WillReturnRows(rows)

	repo := NewPostgresRepository(mock)
	u := &model.User{GoogleID: "gid", Email: "e@x.com", Username: "name", Credit: 10.0}
	if err := repo.Insert(context.Background(), u); err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if u.ID != "user-1" {
		t.Errorf("ID = %q", u.ID)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUsersRepo_Insert_Error(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("INSERT INTO users").WillReturnError(errors.New("dup key"))

	repo := NewPostgresRepository(mock)
	err := repo.Insert(context.Background(), &model.User{Email: "e@x.com"})
	if err == nil {
		t.Error("expected error")
	}
}

func fullUserRows() *pgxmock.Rows {
	now := time.Now()
	avatar := "http://pic"
	return pgxmock.NewRows([]string{
		"id", "google_id", "email", "username", "avatar_url",
		"credit", "member_type", "last_daily_credit_at", "created_at", "updated_at",
	}).AddRow("user-1", "gid", "e@x.com", "name", &avatar, 10.0, "free", (*time.Time)(nil), now, now)
}

func TestUsersRepo_FindByGoogleID_Found(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("SELECT id, google_id").
		WithArgs("gid").
		WillReturnRows(fullUserRows())

	repo := NewPostgresRepository(mock)
	u, err := repo.FindByGoogleID(context.Background(), "gid")
	if err != nil {
		t.Fatalf("FindByGoogleID: %v", err)
	}
	if u.ID != "user-1" || u.Email != "e@x.com" {
		t.Errorf("user = %+v", u)
	}
}

func TestUsersRepo_FindByGoogleID_NotFound(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("SELECT id, google_id").
		WithArgs("gid").
		WillReturnError(pgx.ErrNoRows)

	repo := NewPostgresRepository(mock)
	_, err := repo.FindByGoogleID(context.Background(), "gid")
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("want ErrUserNotFound, got %v", err)
	}
}

func TestUsersRepo_FindByID_Found(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("SELECT id, google_id").
		WithArgs("user-1").
		WillReturnRows(fullUserRows())

	repo := NewPostgresRepository(mock)
	u, err := repo.FindByID(context.Background(), "user-1")
	if err != nil || u.ID != "user-1" {
		t.Errorf("u=%+v err=%v", u, err)
	}
}

func TestUsersRepo_FindByID_NotFound(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("SELECT id, google_id").
		WithArgs("nope").
		WillReturnError(pgx.ErrNoRows)

	repo := NewPostgresRepository(mock)
	if _, err := repo.FindByID(context.Background(), "nope"); !errors.Is(err, ErrUserNotFound) {
		t.Errorf("want ErrUserNotFound, got %v", err)
	}
}

func TestUsersRepo_DeductCredit_Success(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("UPDATE users SET credit").
		WithArgs(1.0, "user-1").
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("user-1"))

	repo := NewPostgresRepository(mock)
	if err := repo.DeductCredit(context.Background(), "user-1", 1.0); err != nil {
		t.Fatalf("DeductCredit: %v", err)
	}
}

func TestUsersRepo_DeductCredit_Insufficient(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()
	mock.ExpectQuery("UPDATE users SET credit").
		WithArgs(9999.0, "user-1").
		WillReturnError(pgx.ErrNoRows)

	repo := NewPostgresRepository(mock)
	if err := repo.DeductCredit(context.Background(), "user-1", 9999.0); err == nil {
		t.Error("expected error for insufficient credit")
	}
}
