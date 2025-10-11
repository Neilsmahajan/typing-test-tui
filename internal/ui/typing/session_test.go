package typing

import (
	"math"
	"testing"
	"time"
)

func TestSessionLifecycle(t *testing.T) {
	session := NewSession()
	if session.Started() {
		t.Fatalf("expected new session to be not started")
	}

	start := time.Now()
	session.Start(start)
	if !session.Started() {
		t.Fatalf("expected session to be started")
	}

	halfway := start.Add(30 * time.Second)
	wpm := session.CurrentWPM(halfway, "hello world")
	if wpm <= 0 {
		t.Fatalf("expected positive WPM, got %f", wpm)
	}

	finish := start.Add(time.Minute)
	final := session.Finish(finish, "hello world")
	if math.Abs(final-2) > 0.0001 {
		t.Fatalf("expected finish WPM to be close to 2, got %f", final)
	}
	if !session.Finished() {
		t.Fatalf("expected session to be finished")
	}
	if session.Elapsed(finish) != time.Minute {
		t.Fatalf("expected elapsed minute, got %s", session.Elapsed(finish))
	}
}

func TestSessionReset(t *testing.T) {
	session := NewSession()
	now := time.Now()
	session.Start(now)
	session.Finish(now.Add(time.Second), "test")

	session.Reset()
	if session.Started() || session.Finished() {
		t.Fatalf("expected session reset to clear state")
	}
	if !session.StartTime().IsZero() || !session.EndTime().IsZero() {
		t.Fatalf("expected reset to clear timestamps")
	}
	if session.WPM() != 0 {
		t.Fatalf("expected reset to clear wpm")
	}
}

func TestSessionFinishZeroDuration(t *testing.T) {
	session := NewSession()
	now := time.Now()
	session.Start(now)
	wpm := session.Finish(now, "word")
	if wpm != 0 {
		t.Fatalf("expected zero wpm for zero duration finish, got %f", wpm)
	}
}
