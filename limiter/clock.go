package limiter

import "time"

type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// realClock production environment uses real time
type realClock struct{}

func (realClock) Now() time.Time        { return time.Now() }
func (realClock) Sleep(d time.Duration) { time.Sleep(d) }

// mockClock is a controllable time implementation for testing (not thread-safe, used in single-threaded tests)
type mockClock struct {
	now time.Time
}

func NewMockClock(start time.Time) *mockClock { return &mockClock{now: start} }

func (m *mockClock) Now() time.Time          { return m.now }
func (m *mockClock) Sleep(d time.Duration)   { m.Advance(d) }
func (m *mockClock) Advance(d time.Duration) { m.now = m.now.Add(d) }
