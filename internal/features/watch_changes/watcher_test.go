package watch_changes_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MasayukiFukada/PjDoc/internal/features/watch_changes"
)

func TestBroadcaster(t *testing.T) {
	b := watch_changes.NewBroadcaster()
	ch := make(chan string, 1)

	b.Register(ch)
	b.NotifyReload()

	select {
	case msg := <-ch:
		if msg != "reload" {
			t.Errorf("Expected 'reload', got %q", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for reload notification")
	}

	b.Unregister(ch)
}

func TestHandler(t *testing.T) {
	b := watch_changes.NewBroadcaster()
	handler := watch_changes.NewHandler(b)

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/api/events", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	go func() {
		time.Sleep(50 * time.Millisecond)
		b.NotifyReload()
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	handler.ServeHTTP(rec, req)

	output := rec.Body.String()
	expected := "data: reload\n\n"
	if output != expected {
		t.Errorf("Expected body %q, got %q", expected, output)
	}
}
