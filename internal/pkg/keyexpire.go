package pkg

import (
	"time"
)

type KeyExpirer struct {
	redis *Redis
	done  chan bool
}

func NewKeyExpirer(r *Redis) *KeyExpirer {
	return &KeyExpirer{
		redis: r,
		done:  make(chan bool),
	}
}

// Start starts the Expirer.
//
// tick - How fast is the cleaner triggered.
//
// randomKeys - Amount of random expiring keys to get checked.
//
// againPercentage - If more than x% of keys were expired, start again in same tick.
func (e *KeyExpirer) Start(tick time.Duration) {
	ticker := time.NewTicker(tick)
	for {
		select {
		case <-ticker.C:
			{
				for _, db := range e.redis.RedisDbs() {
					e.redis.mu.Lock()
					db.DeleteExpiredKeys()
					e.redis.mu.Lock()
				}
			}
		case <-e.done:
			ticker.Stop()
			return
		}
	}
}

func (ke *KeyExpirer) Stop() {
	if ke.done != nil {
		ke.done <- true
		close(ke.done)
	}
}
