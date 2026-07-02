# เครียดโว้ยยยย (Khriat-woeyyyyy)

> "แปลงคำบ่นเป็นภาษาคนมีการศึกษา" — Turn your workplace rage into polished corporate speech, powered by AI.

Paste whatever you actually want to say. Pick who you're sending it to and how formal you want to sound. The AI rewrites it for you — from aristocratic politeness to chaotic unhinged energy.

---

## What it does

You type the raw, frustrated version of what you want to say. The app rewrites it through one of five personality modes targeted at whoever you're sending it to.

**Recipients**
- หัวหน้า — Boss
- ลูกค้า — Client
- เพื่อน — Friend

**Personality levels**

| Level | Name | Vibe |
|-------|------|------|
| 1 | สวมวิญญาณผู้ดี | Dripping with aristocratic politeness — but the disdain is visible |
| 2 | พูดดีด้วยละนะ | Passive-aggressive admin who has answered this question 100 times today |
| 3 | มนุษย์ปกติ | Normal, friendly, professional — nothing weird |
| 4 | นึกว่าสนิท | Acts like you've been best friends since birth |
| 5 | ตัวมัมซัมซุง | Chaotic, typos, unhinged emotional energy |

Output is in Thai by default. The backend also supports English output (`lang: "en"` in the API body).

---

## How credits work

| State | Credits | Resets |
|-------|---------|--------|
| Guest (no login) | 5 free translations | Per browser session (Redis TTL) |
| Logged in via Google | More credits | Stored in your account (PostgreSQL) |

Credits are deducted only on a successful AI response. Rate limiting is enforced at 10 translate requests per minute per IP.

**Topping up** — logged-in users can buy more credits via PromptPay QR (Omise). 1 THB paid = 1 credit granted, credited automatically once the payment settles. Guests must log in first (credits are tied to an account, not a browser session).

1. Click **เติมเครดิต** (Top up) in the navbar, or on the home page once credits hit 0
2. Enter an amount (20–5,000 THB) and generate a QR code
3. Scan and pay with any Thai banking app that supports PromptPay
4. The frontend polls payment status every 3 seconds; credit updates automatically on success

---

## Tech stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 15, React 19, TypeScript, Tailwind CSS |
| Backend | Go 1.25, Fiber v2 |
| AI | DeepSeek API (OpenAI-compatible) |
| Auth | Google Identity Services → backend-issued JWT |
| Payments | Omise (PromptPay QR), HMAC-SHA256 + IP-allowlisted webhook |
| Database | PostgreSQL (pgx/pgxpool, raw SQL, golang-migrate) |
| Sessions | Redis (guest credit tracking) |
| Deploy — backend | Render (Docker) |
| Deploy — frontend | Vercel |

---

## Local development

### Prerequisites

- Go 1.25+
- Node.js 20+
- Docker + Docker Compose
- A DeepSeek API key ([platform.deepseek.com](https://platform.deepseek.com))
- A Google OAuth Client ID ([console.cloud.google.com](https://console.cloud.google.com))

### 1. Create the env file

Create `.env` in the project root:

```env
DATABASE_URL=postgres://postgres:secret@127.0.0.1:5432/myapp_db?sslmode=disable
PORT=8080
REDIS_URL=localhost:6379
FRONTEND_ORIGIN=http://localhost:3000
JWT_SECRET=any-long-random-string-you-choose
GOOGLE_CLIENT_ID=your-google-oauth-client-id.apps.googleusercontent.com
DEEPSEEK_API_KEY=sk-your-deepseek-key
AUTO_MIGRATE=true
APP_ENV=development
OMISE_SECRET_KEY=skey_test_your-omise-secret-key
# Omise's published webhook source IPs — https://docs.omise.co/api-ips
OMISE_WEBHOOK_ALLOWED_IPS=54.169.118.227,52.74.199.175,18.139.13.19
# Omise dashboard → Webhooks → "Roll secret" (base64-encoded value)
OMISE_WEBHOOK_SECRET=your-base64-webhook-secret
```

Create `apps/web/.env.local`:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_GOOGLE_CLIENT_ID=your-google-oauth-client-id.apps.googleusercontent.com
```

### 2. Start Postgres and Redis

```bash
docker compose up -d postgres redis
```

### 3. Run the backend

```bash
cd apps/api
go run ./cmd/api
```

API is at `http://localhost:8080`. Migrations run automatically on first start (`AUTO_MIGRATE=true`).

### 4. Run the frontend

```bash
cd apps/web
npm install
npm run dev
```

Frontend is at `http://localhost:3000`.

### Google Sign-In (local)

Add both `http://localhost:3000` and `http://127.0.0.1:3000` as **Authorized JavaScript Origins** in your Google Cloud Console OAuth client. The frontend must run on port 3000 — changing the port will cause Google Sign-In to fail with `unregistered_origin`.

---

## Running tests

**Backend** (85.9% line coverage):

```bash
cd apps/api
go test ./...
```

**Frontend** (93.2% line coverage):

```bash
cd apps/web
npm test
```

---

## Production deployment

### Backend → Render

`render.yaml` is pre-configured. Connect the GitHub repo in the Render dashboard (New → Web Service → select this repo). Set these environment variables in the Render dashboard — **do not put real values in code or `render.yaml`**:

| Variable | Where to get it |
|----------|----------------|
| `DATABASE_URL` | Supabase → Settings → Database → Connection string |
| `REDIS_URL` | Upstash → your Redis instance → `rediss://` URL |
| `JWT_SECRET` | Generate with `openssl rand -hex 32` |
| `GOOGLE_CLIENT_ID` | Google Cloud Console → your OAuth client |
| `DEEPSEEK_API_KEY` | DeepSeek platform dashboard |
| `FRONTEND_ORIGIN` | Your Vercel URL (set after frontend is deployed) |
| `OMISE_SECRET_KEY` | Omise dashboard → Keys (`skey_live_...` / `skey_test_...`) |
| `OMISE_WEBHOOK_ALLOWED_IPS` | `54.169.118.227,52.74.199.175,18.139.13.19` — Omise's published webhook IPs ([docs.omise.co/api-ips](https://docs.omise.co/api-ips)); re-check that page before going live, IPs can change |
| `OMISE_WEBHOOK_SECRET` | Omise dashboard → Webhooks → "Roll secret" (base64-encoded value) |

### Frontend → Vercel

`vercel.json` is pre-configured (root directory = `apps/web`). Import the GitHub repo at [vercel.com/new](https://vercel.com/new). Set these environment variables in the Vercel dashboard:

| Variable | Value |
|----------|-------|
| `NEXT_PUBLIC_API_BASE_URL` | Your Render service URL |
| `NEXT_PUBLIC_GOOGLE_CLIENT_ID` | Your Google OAuth Client ID |

### After both are live

Add your production Vercel URL as an **Authorized JavaScript Origin** in Google Cloud Console → APIs & Services → Credentials → your OAuth client. Without this, Google Sign-In will fail in production with `unregistered_origin`.

---

## API reference

### `GET /health`
Returns `{"status":"ok"}`. No auth required.

### `GET /guest/status`
Returns current credit balance. Works for both guests (cookie session) and logged-in users (JWT).

### `POST /translate`
Rate limited: 10 requests per minute per IP. Max request body: 64 KB. Max text length: 3,000 characters.

Request body:
```json
{
  "text": "ทำไมส่งงานช้าจัง",
  "target": "boss",
  "level": 3,
  "lang": "th"
}
```

- `target`: `"boss"` | `"client"` | `"friend"`
- `level`: `1`–`5` (default `3`)
- `lang`: `"th"` (default) | `"en"`

Response:
```json
{
  "result": "ขออภัยที่งานล่าช้า จะเร่งดำเนินการให้เสร็จโดยเร็วที่สุดครับ",
  "level": 3,
  "target": "boss"
}
```

Error codes:
- `400` — missing/invalid body, or text over 3,000 characters
- `401` — session expired
- `402` — credit exhausted (guest: log in to get more; user: insufficient credit)
- `429` — rate limit hit

### `POST /api/v1/auth/google`
Exchange a Google ID token for a backend JWT. Rate limited: 20 requests per minute per IP.

### `POST /api/v1/payments/create`
Requires auth (`Authorization: Bearer <jwt>`). Creates an Omise PromptPay charge and a pending payment record.

Request body:
```json
{ "amount": 50 }
```
- `amount`: THB, must be between 20 and 5,000

Response (`201`):
```json
{
  "paymentId": "3f2b...",
  "status": "pending",
  "amount": 50,
  "currency": "THB",
  "qrCodeUri": "https://..."
}
```

### `GET /api/v1/payments/:id/status`
Requires auth. Returns the payment status, scoped to the requesting user (`404` if it belongs to someone else).

```json
{ "paymentId": "3f2b...", "status": "pending", "amount": 50 }
```
`status` is one of `pending` | `success` | `failed`.

### `POST /webhooks/omise`
Called by Omise, not the frontend. Verifies the `Omise-Signature-Timestamp` / `Omise-Signature` headers (HMAC-SHA256, base64-decoded webhook secret — see [docs.omise.co/api-webhooks](https://docs.omise.co/api-webhooks#protecting-your-endpoints)) and restricts the source IP to Omise's published webhook servers. On a successful charge, credits the user's balance and records a `credit_ledger` entry in a single transaction. Idempotent — replays of the same event are skipped once processed.

---

## Project structure

```
.
├── apps/
│   ├── api/                   # Go backend
│   │   ├── cmd/api/           # Entry point (main.go)
│   │   ├── internal/
│   │   │   ├── config/        # Env loading
│   │   │   ├── database/      # Postgres + Redis + migrations
│   │   │   ├── handler/       # HTTP handlers
│   │   │   ├── middleware/    # Guest session, JWT auth, Omise IP allowlist
│   │   │   ├── model/         # Domain structs
│   │   │   ├── repository/    # Data access layer
│   │   │   ├── routes/        # Route registration
│   │   │   └── service/       # Business logic, AI calls, Omise client
│   │   └── Dockerfile
│   └── web/                   # Next.js frontend
│       ├── app/               # App Router pages + contexts (incl. payment/)
│       ├── components/        # Navbar, PillBar, etc.
│       └── fonts/
├── postman/                   # Postman collection for manual API testing
├── docker-compose.yml         # Local Postgres + Redis
├── render.yaml                # Render deployment config
└── vercel.json                # Vercel deployment config
```
