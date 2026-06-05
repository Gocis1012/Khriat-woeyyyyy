# Corporate Translator

A full-stack MVP for translating emotionally charged workplace messages into polished corporate communication.

## Stack

- Frontend: Next.js, React, TypeScript
- Backend: Go Fiber
- Database: PostgreSQL in Docker
- SQL: raw SQL with `pgxpool`
- Auth: Google OAuth ID token -> backend-issued JWT
- AI: OpenAI Chat Completions

## Local Setup

1. Copy `.env.example` into the backend and frontend env files:

```bash
cp .env.example apps/api/.env
cp .env.example apps/web/.env.local
```

2. Start Postgres:

```bash
docker compose up -d postgres
```

3. Run the API:

```bash
cd apps/api
go run ./cmd/api
```

4. Run the web app:

```bash
cd apps/web
npm install
npm run dev
```

The frontend runs at `http://localhost:3000`, the API runs at `http://localhost:8080`, and Postgres is exposed on `127.0.0.1:5433`.

## Notes

- Google OAuth is the only login method.
- The backend verifies Google ID tokens and then issues its own JWT.
- Migrations are `.up.sql` and `.down.sql` files under `apps/api/internal/migrate/migrations`.
- If `OPENAI_API_KEY` is empty, the API returns a deterministic demo translation so local development can run without external secrets.
