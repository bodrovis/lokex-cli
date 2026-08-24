package commandctx

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"
)

func TestNewCommandContext_NonPositiveTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "zero",
			timeout: 0,
		},
		{
			name:    "negative",
			timeout: -time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := NewCommandContext(tt.timeout)

			if ctx == nil {
				t.Fatal("NewCommandContext() context = nil")
			}

			if _, ok := ctx.Deadline(); ok {
				t.Fatal("NewCommandContext() set a deadline, want none")
			}

			if ctx.Done() != nil {
				t.Fatal("NewCommandContext() Done() != nil, want nil")
			}

			cancel()

			if err := ctx.Err(); err != nil {
				t.Fatalf("context error after cancel = %v, want nil", err)
			}
		})
	}
}

func TestNewCommandContext_PositiveTimeoutSetsDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const timeout = 200 * time.Millisecond

		start := time.Now()

		ctx, cancel := NewCommandContext(timeout)
		defer cancel()

		deadline, ok := ctx.Deadline()
		if !ok {
			t.Fatal("NewCommandContext() did not set deadline")
		}

		want := start.Add(timeout)
		if !deadline.Equal(want) {
			t.Fatalf("deadline = %v, want %v", deadline, want)
		}

		if err := ctx.Err(); err != nil {
			t.Fatalf("context error = %v, want nil", err)
		}
	})
}

func TestNewCommandContext_PositiveTimeoutCanBeCanceled(t *testing.T) {
	ctx, cancel := NewCommandContext(time.Hour)
	cancel()

	select {
	case <-ctx.Done():
	default:
		t.Fatal("context was not canceled")
	}

	if !errors.Is(ctx.Err(), context.Canceled) {
		t.Fatalf("context error = %v, want %v", ctx.Err(), context.Canceled)
	}
}

func TestNewCommandContext_PositiveTimeoutExpires(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const timeout = time.Hour

		ctx, cancel := NewCommandContext(timeout)
		defer cancel()

		synctest.Sleep(timeout)

		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Fatalf(
				"context error = %v, want %v",
				ctx.Err(),
				context.DeadlineExceeded,
			)
		}
	})
}
