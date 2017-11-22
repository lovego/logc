package elastic_search

import (
	"time"
)

type Timer struct {
	duration time.Duration
}

func (t *Timer) Sleep() {
	const max = 10 * time.Minute
	if t.duration <= 0 {
		t.duration = time.Second
	} else if t.duration < max {
		if t.duration *= 2; t.duration > max {
			t.duration = max
		}
	}
	time.Sleep(t.duration)
}
