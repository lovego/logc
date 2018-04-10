package elasticsearch

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	var expect = time.Second
	timer := Timer{}
	for i := 0; i < 2; i++ {
		timer.Sleep()
		if timer.duration != expect {
			t.Fatalf("expect: %v, got: %v.", expect, timer.duration)
		}
		if expect *= 2; expect > 10*time.Minute {
			expect = 10 * time.Minute
		}
	}
}
