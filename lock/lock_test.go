package lock

import (
	"context"
	"testing"
	"time"
)

func TestTrySuccess(t *testing.T) {
	lock := New()
	if !tryWithTimeout(lock, "1ms") {
		t.Fatalf("Try should succeed")
	}
}

func TestTryFailure(t *testing.T) {
	lock := New()
	tryWithTimeout(lock, "1ms")
	if tryWithTimeout(lock, "1ms") {
		t.Fatalf("Try should fail")
	}
}

func TestTryAfterUnlock(t *testing.T) {
	lock := New()
	if !tryWithTimeout(lock, "1ms") {
		t.Fatalf("Try should succeed")
	}
	if tryWithTimeout(lock, "1ms") {
		t.Fatalf("Try should fail")
	}
	lock.Unlock()
	if !tryWithTimeout(lock, "1ms") {
		t.Fatalf("Try should succeed")
	}
}

func tryWithTimeout(lock Lock, timeout string) bool {
	duration, _ := time.ParseDuration(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return lock.Try(ctx)
}
