package typing

import (
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	DefaultBoxWidth     = 60
	BoxHorizontalMargin = 4

	averageWordLength = 5
)

type Session struct {
	started  bool
	start    time.Time
	finished bool
	end      time.Time
	wpm      float64
}

func NewSession() Session {
	return Session{}
}

func (s *Session) Reset() {
	s.started = false
	s.finished = false
	s.start = time.Time{}
	s.end = time.Time{}
	s.wpm = 0
}

func (s *Session) Start(now time.Time) {
	if s.started {
		return
	}
	s.started = true
	s.start = now
}

func (s *Session) Finish(now time.Time, target string) float64 {
	if s.finished {
		return s.wpm
	}
	s.finished = true
	s.end = now
	elapsedMinutes := s.Elapsed(now).Minutes()
	if elapsedMinutes <= 0 {
		s.wpm = 0
		return s.wpm
	}
	words := float64(len(strings.Fields(target)))
	if words <= 0 {
		runes := utf8.RuneCountInString(target)
		if runes > 0 {
			words = float64(runes) / averageWordLength
		}
	}
	if words <= 0 {
		s.wpm = 0
		return s.wpm
	}
	s.wpm = words / elapsedMinutes
	if math.IsNaN(s.wpm) || math.IsInf(s.wpm, 0) {
		s.wpm = 0
	}
	return s.wpm
}

func (s *Session) CurrentWPM(now time.Time, typed string) float64 {
	if s.finished {
		return s.wpm
	}
	if !s.started {
		return 0
	}
	elapsed := now.Sub(s.start).Minutes()
	if elapsed <= 0 {
		return 0
	}
	words := float64(len(strings.Fields(typed)))
	if words == 0 {
		words = float64(utf8.RuneCountInString(typed)) / averageWordLength
	}
	if words <= 0 {
		return 0
	}
	value := words / elapsed
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}
	return value
}

func (s *Session) Elapsed(now time.Time) time.Duration {
	if !s.started {
		return 0
	}
	if s.finished {
		return s.end.Sub(s.start)
	}
	return now.Sub(s.start)
}

func (s *Session) Started() bool {
	return s.started
}

func (s *Session) Finished() bool {
	return s.finished
}

func (s *Session) StartTime() time.Time {
	return s.start
}

func (s *Session) EndTime() time.Time {
	return s.end
}

func (s *Session) WPM() float64 {
	return s.wpm
}
