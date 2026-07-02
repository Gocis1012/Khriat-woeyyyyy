package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"testing"
	"time"
)

// sign computes the same "<timestamp>.<rawBody>" HMAC-SHA256 hex signature
// the real Omise webhook sender would produce, for use as test fixtures.
func sign(t *testing.T, base64Secret, timestamp, rawBody string) string {
	t.Helper()
	secret, err := base64.StdEncoding.DecodeString(base64Secret)
	if err != nil {
		t.Fatalf("bad test secret: %v", err)
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(timestamp + "." + rawBody))
	return hex.EncodeToString(mac.Sum(nil))
}

func TestPaymentService_VerifyWebhookSignature_Valid(t *testing.T) {
	const secret = "d2hzZWNfdGVzdA==" // base64("whsec_test")
	svc := NewPaymentService(nil, nil, secret)

	body := `{"data":{"id":"chrg_test","status":"successful"}}`
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := sign(t, secret, ts, body)

	if !svc.VerifyWebhookSignature([]byte(body), ts, sig) {
		t.Error("expected valid signature to verify")
	}
}

func TestPaymentService_VerifyWebhookSignature_WrongSignature(t *testing.T) {
	const secret = "d2hzZWNfdGVzdA=="
	svc := NewPaymentService(nil, nil, secret)

	body := `{"data":{"id":"chrg_test","status":"successful"}}`
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	if svc.VerifyWebhookSignature([]byte(body), ts, "0000000000000000000000000000000000000000000000000000000000000000") {
		t.Error("expected mismatched signature to fail")
	}
}

func TestPaymentService_VerifyWebhookSignature_WrongSecret(t *testing.T) {
	svc := NewPaymentService(nil, nil, "d2hzZWNfdGVzdA==")

	body := `{"data":{"id":"chrg_test","status":"successful"}}`
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	// Signed with a different secret than the service holds.
	sig := sign(t, "b3RoZXJfc2VjcmV0", ts, body)

	if svc.VerifyWebhookSignature([]byte(body), ts, sig) {
		t.Error("expected signature signed with a different secret to fail")
	}
}

func TestPaymentService_VerifyWebhookSignature_TamperedBody(t *testing.T) {
	const secret = "d2hzZWNfdGVzdA=="
	svc := NewPaymentService(nil, nil, secret)

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig := sign(t, secret, ts, `{"data":{"id":"chrg_test","status":"successful"}}`)

	tampered := `{"data":{"id":"chrg_test","status":"failed"}}`
	if svc.VerifyWebhookSignature([]byte(tampered), ts, sig) {
		t.Error("expected signature over the original body to fail against a tampered body")
	}
}

func TestPaymentService_VerifyWebhookSignature_StaleTimestamp(t *testing.T) {
	const secret = "d2hzZWNfdGVzdA=="
	svc := NewPaymentService(nil, nil, secret)

	body := `{"data":{"id":"chrg_test","status":"successful"}}`
	old := strconv.FormatInt(time.Now().Add(-10*time.Minute).Unix(), 10)
	sig := sign(t, secret, old, body)

	if svc.VerifyWebhookSignature([]byte(body), old, sig) {
		t.Error("expected stale (replayed) timestamp to be rejected")
	}
}

func TestPaymentService_VerifyWebhookSignature_FutureTimestamp(t *testing.T) {
	const secret = "d2hzZWNfdGVzdA=="
	svc := NewPaymentService(nil, nil, secret)

	body := `{"data":{"id":"chrg_test","status":"successful"}}`
	future := strconv.FormatInt(time.Now().Add(10*time.Minute).Unix(), 10)
	sig := sign(t, secret, future, body)

	if svc.VerifyWebhookSignature([]byte(body), future, sig) {
		t.Error("expected far-future timestamp to be rejected")
	}
}

func TestPaymentService_VerifyWebhookSignature_MalformedTimestamp(t *testing.T) {
	svc := NewPaymentService(nil, nil, "d2hzZWNfdGVzdA==")
	if svc.VerifyWebhookSignature([]byte("{}"), "not-a-number", "deadbeef") {
		t.Error("expected malformed timestamp to be rejected")
	}
}

func TestPaymentService_VerifyWebhookSignature_RotationSecondSignatureMatches(t *testing.T) {
	const secret = "d2hzZWNfdGVzdA=="
	svc := NewPaymentService(nil, nil, secret)

	body := `{"data":{"id":"chrg_test","status":"successful"}}`
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	realSig := sign(t, secret, ts, body)
	header := "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef," + realSig

	if !svc.VerifyWebhookSignature([]byte(body), ts, header) {
		t.Error("expected the second comma-separated signature to be accepted during rotation")
	}
}

func TestPaymentService_VerifyWebhookSignature_EmptyInputs(t *testing.T) {
	svc := NewPaymentService(nil, nil, "d2hzZWNfdGVzdA==")

	if svc.VerifyWebhookSignature([]byte("{}"), "", "sig") {
		t.Error("expected empty timestamp header to be rejected")
	}
	if svc.VerifyWebhookSignature([]byte("{}"), "123", "") {
		t.Error("expected empty signature header to be rejected")
	}

	svcNoSecret := NewPaymentService(nil, nil, "")
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	if svcNoSecret.VerifyWebhookSignature([]byte("{}"), ts, "sig") {
		t.Error("expected verification to fail when no webhook secret is configured")
	}
}

func TestPaymentService_VerifyWebhookSignature_InvalidBase64Secret(t *testing.T) {
	svc := NewPaymentService(nil, nil, "not-valid-base64!!!")
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	if svc.VerifyWebhookSignature([]byte("{}"), ts, "sig") {
		t.Error("expected invalid base64 secret to fail closed")
	}
}
