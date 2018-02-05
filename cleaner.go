package ratelimit

import (
	"time"
)

func RunLimitCleaner(l *Limit) {
	go func(l *Limit) {
		for {
			time.Sleep(time.Second)
			now := time.Now().Add(-l.Per)
			Mutex.Lock()
			for k, r := range l.Rates {
				if r.ExpiredAt.Before(now) {
					l.Rates[k] = createKey()
				}
			}
			Mutex.Unlock()
		}
	}(l)
}
