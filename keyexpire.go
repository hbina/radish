package redis

import (
	"math"
	"time"
)

type KeyExpirer interface {
	Start(tick time.Duration)
	Stop()
}

var _ KeyExpirer = (*Expirer)(nil)

type Expirer struct {
	redis *Redis

	done chan bool
}

func NewKeyExpirer(r *Redis) *Expirer {
	return &Expirer{
		redis: r,
		done:  make(chan bool, math.MaxInt32),
	}
}

// Start starts the Expirer.
//
// tick - How fast is the cleaner triggered.
//
// randomKeys - Amount of random expiring keys to get checked.
//
// againPercentage - If more than x% of keys were expired, start again in same tick.
func (e *Expirer) Start(tick time.Duration) {
	ticker := time.NewTicker(tick)
	for {
		select {
		case <-ticker.C:
			e.cleanupExpiredKeys()
		case <-e.done:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the
func (e *Expirer) Stop() {
	if e.done != nil {
		e.done <- true
		close(e.done)
	}
}

func (e *Expirer) cleanupExpiredKeys() {
	var count int = 0
	e.Redis().mu.Lock()
	defer e.Redis().mu.Unlock()
	for _, db := range e.Redis().RedisDbs() {
		count += db.DeleteExpiredKeys()
	}
}

// Redis gets the redis instance.
func (e *Expirer) Redis() *Redis {
	return e.redis
}
