package database

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestNewRedisClient_HostPort(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client, err := NewRedisClient(mr.Addr())
	if err != nil {
		t.Fatalf("NewRedisClient(host:port): %v", err)
	}
	defer client.Close()
}

func TestNewRedisClient_RedisURL(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	client, err := NewRedisClient("redis://" + mr.Addr())
	if err != nil {
		t.Fatalf("NewRedisClient(redis://...): %v", err)
	}
	defer client.Close()
}

func TestNewRedisClient_InvalidURL(t *testing.T) {
	if _, err := NewRedisClient("redis://%zz-not-valid"); err == nil {
		t.Error("expected error for malformed redis URL")
	}
}
