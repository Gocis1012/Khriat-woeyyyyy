package handler

import (
	"context"

	"corporate-translator-api/internal/model"

	"github.com/gofiber/fiber/v2"
)

// ── Fakes ─────────────────────────────────────────────────────────────────────

type fakeUserSvc struct {
	insert         func(ctx context.Context, u *model.User) error
	findByGoogleID func(ctx context.Context, gid string) (*model.User, error)
	findByID       func(ctx context.Context, id string) (*model.User, error)
	deduct         func(ctx context.Context, id string, amt float64) error
}

func (f *fakeUserSvc) Insert(ctx context.Context, u *model.User) error { return f.insert(ctx, u) }
func (f *fakeUserSvc) FindByGoogleID(ctx context.Context, gid string) (*model.User, error) {
	return f.findByGoogleID(ctx, gid)
}
func (f *fakeUserSvc) FindByID(ctx context.Context, id string) (*model.User, error) {
	return f.findByID(ctx, id)
}
func (f *fakeUserSvc) DeductCredit(ctx context.Context, id string, amt float64) error {
	return f.deduct(ctx, id, amt)
}

type fakeGuestSvc struct {
	getStatus func(ctx context.Context, id string) (*model.Guest, error)
	useCredit func(ctx context.Context, id string, fn func() error) error
	del       func(ctx context.Context, id string) error
}

func (f *fakeGuestSvc) GetStatus(ctx context.Context, id string) (*model.Guest, error) {
	return f.getStatus(ctx, id)
}
func (f *fakeGuestSvc) UseCredit(ctx context.Context, id string, fn func() error) error {
	return f.useCredit(ctx, id, fn)
}
func (f *fakeGuestSvc) DeleteSession(ctx context.Context, id string) error { return f.del(ctx, id) }

type fakeTranslator struct {
	purify func(ctx context.Context, text, target string, level int) (string, error)
}

func (f *fakeTranslator) PurifyText(ctx context.Context, text, target string, level int) (string, error) {
	return f.purify(ctx, text, target, level)
}

type fakeAuthSvc struct {
	validateGoogle func(ctx context.Context, idToken string) (*model.GoogleProfile, error)
	generateJWT    func(userID, email string) (string, error)
}

func (f *fakeAuthSvc) ValidateGoogleToken(ctx context.Context, idToken string) (*model.GoogleProfile, error) {
	return f.validateGoogle(ctx, idToken)
}
func (f *fakeAuthSvc) GenerateJWT(userID, email string) (string, error) {
	return f.generateJWT(userID, email)
}

// ── Helper: build a Fiber app that injects Locals from test headers ──────────

func testApp(register func(app *fiber.App)) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		if v := c.Get("X-User"); v != "" {
			c.Locals("user_id", v)
		}
		if v := c.Get("X-Guest"); v != "" {
			c.Locals("guest_id", v)
		}
		return c.Next()
	})
	register(app)
	return app
}
