package timers

import (
	"sync/atomic"
	"time"
)

var lastID int64

func NextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

type TimerMsg struct {
	ID      int
	Timeout bool
	Tag     time.Time
}
