package commandctx

import (
	"testing"
	"time"
)

func TestNewCommandContext_ZeroTimeoutReturnsBackgroundLikeContext(t *testing.T) {
	ctx, cancel := NewCommandContext(0)
	defer cancel()

	if ctx == nil {
		t.Fatal("expected context, got nil")
	}

	if deadline, ok := ctx.Deadline(); ok {
		t.Fatalf("expected no deadline, got %v", deadline)
	}

	if done := ctx.Done(); done != nil {
		t.Fatal("expected nil Done channel for background context")
	}

	if err := ctx.Err(); err != nil {
		t.Fatalf("expected nil error before cancel, got %v", err)
	}

	cancel()

	if err := ctx.Err(); err != nil {
		t.Fatalf("expected nil error after noop cancel, got %v", err)
	}
}

func TestNewCommandContext_NegativeTimeoutReturnsBackgroundLikeContext(t *testing.T) {
	ctx, cancel := NewCommandContext(-time.Second)
	defer cancel()

	if ctx == nil {
		t.Fatal("expected context, got nil")
	}

	if deadline, ok := ctx.Deadline(); ok {
		t.Fatalf("expected no deadline, got %v", deadline)
	}

	if done := ctx.Done(); done != nil {
		t.Fatal("expected nil Done channel for background context")
	}

	if err := ctx.Err(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNewCommandContext_PositiveTimeoutSetsDeadline(t *testing.T) {
	timeout := 200 * time.Millisecond

	before := time.Now()
	ctx, cancel := NewCommandContext(timeout)
	defer cancel()
	after := time.Now()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}

	minDeadline := before.Add(timeout)
	maxDeadline := after.Add(timeout)

	if deadline.Before(minDeadline) {
		t.Fatalf("deadline is too early: got %v, want >= %v", deadline, minDeadline)
	}

	if deadline.After(maxDeadline) {
		t.Fatalf("deadline is too late: got %v, want <= %v", deadline, maxDeadline)
	}

	if err := ctx.Err(); err != nil {
		t.Fatalf("expected nil error before timeout/cancel, got %v", err)
	}
}

func TestNewCommandContext_PositiveTimeoutCanBeCanceled(t *testing.T) {
	ctx, cancel := NewCommandContext(time.Hour)

	cancel()

	select {
	case <-ctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected context to be canceled")
	}

	if err := ctx.Err(); err == nil {
		t.Fatal("expected context error after cancel")
	}
}

func TestNewCommandContext_PositiveTimeoutExpires(t *testing.T) {
	ctx, cancel := NewCommandContext(10 * time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Fatal("expected context to expire")
	}

	if err := ctx.Err(); err == nil {
		t.Fatal("expected context error after timeout")
	}
}
