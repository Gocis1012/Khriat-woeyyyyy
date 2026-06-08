package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

type fakeCompleter struct {
	resp   openai.ChatCompletionResponse
	err    error
	gotReq openai.ChatCompletionRequest
}

func (f *fakeCompleter) CreateChatCompletion(
	_ context.Context,
	req openai.ChatCompletionRequest,
) (openai.ChatCompletionResponse, error) {
	f.gotReq = req
	return f.resp, f.err
}

func respWith(content string) openai.ChatCompletionResponse {
	return openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{Message: openai.ChatCompletionMessage{Content: content}},
		},
	}
}

func TestNewTranslationService(t *testing.T) {
	svc, err := NewTranslationService("sk-test")
	if err != nil || svc == nil {
		t.Fatalf("svc=%v err=%v", svc, err)
	}
}

func TestPurifyText_Success(t *testing.T) {
	fake := &fakeCompleter{resp: respWith("สุภาพแล้วครับ")}
	svc := newWithClient(fake)

	out, err := svc.PurifyText(context.Background(), "ด่าๆ", "boss", 3, "th")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "สุภาพแล้วครับ" {
		t.Errorf("output = %q", out)
	}
	// The user message should carry the raw input.
	last := fake.gotReq.Messages[len(fake.gotReq.Messages)-1]
	if last.Content != "ด่าๆ" {
		t.Errorf("user message = %q", last.Content)
	}
}

func TestPurifyText_EmptyChoices(t *testing.T) {
	svc := newWithClient(&fakeCompleter{resp: openai.ChatCompletionResponse{}})
	if _, err := svc.PurifyText(context.Background(), "x", "boss", 3, "th"); err == nil {
		t.Error("expected error for empty choices")
	}
}

func TestPurifyText_ClientError(t *testing.T) {
	svc := newWithClient(&fakeCompleter{err: errors.New("api down")})
	if _, err := svc.PurifyText(context.Background(), "x", "boss", 3, "th"); err == nil {
		t.Error("expected error from client")
	}
}

func TestPurifyText_DefaultsLevelAndTarget(t *testing.T) {
	fake := &fakeCompleter{resp: respWith("ok")}
	svc := newWithClient(fake)

	// level 0 and 99 should fall back to 3; empty target -> boss
	for _, lvl := range []int{0, 99} {
		if _, err := svc.PurifyText(context.Background(), "x", "", lvl, "th"); err != nil {
			t.Errorf("level %d: unexpected error %v", lvl, err)
		}
	}
	// System prompt should reference boss context (เจ้านาย)
	sys := fake.gotReq.Messages[0].Content
	if !strings.Contains(sys, "เจ้านาย") {
		t.Errorf("system prompt missing boss context: %q", sys)
	}
}

func TestPurifyText_AllTargetsAndLevels(t *testing.T) {
	fake := &fakeCompleter{resp: respWith("ok")}
	svc := newWithClient(fake)
	for _, target := range []string{"boss", "client", "teacher", "friend", "unknown"} {
		for lvl := 1; lvl <= 5; lvl++ {
			for _, lang := range []string{"th", "en", ""} {
				if _, err := svc.PurifyText(context.Background(), "x", target, lvl, lang); err != nil {
					t.Errorf("target=%s level=%d lang=%s: %v", target, lvl, lang, err)
				}
			}
		}
	}
}
