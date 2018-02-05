package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type Limit struct {
	MaxRequests int
	Per         time.Duration
	Block       time.Duration
	Blocker     Blocker
	MaxSpam     int
	Spammer     Spammer
	Rates       map[string]*RateLimit
}
type RateLimit struct {
	ExpiredAt time.Time
	Hits      int
}

var (
	Mutex sync.Mutex
)

func CreateLimit(key string) Limit {
	op, err := parse(key)
	if err != nil {
		panic(fmt.Sprintf("Faild to parse %s : %q", key, err))
	}

	limits := make(map[string]*RateLimit)
	l := Limit{
		MaxRequests: op.Max,
		Per:         op.Per,
		Block:       op.Block,
		MaxSpam:     op.MaxToSpam,
		Rates:       limits,
	}
	RunLimitCleaner(&l)

	if l.MaxSpam != 0 {
		l.Spammer = CreateSpammer()
	}
	if l.Block != 0 {
		l.Blocker = CreateBlocker()
	}

	return l
}

func createKey() *RateLimit {
	return &RateLimit{
		ExpiredAt: time.Now(),
	}
}

func (l *Limit) Hit(key string) error {
	Mutex.Lock()
	k, ok := l.Rates[key]
	if !ok {
		l.Rates[key] = createKey()
		k = l.Rates[key]
	}
	if k.Hits >= l.MaxRequests {
		if l.Spammer.Values != nil {
			l.Spammer.Increase(key)
		}
		if l.Spammer.Values != nil && l.Blocker.Values != nil {
			if l.Spammer.Values[key].Hits >= l.MaxSpam {
				l.Blocker.AddIfNotExists(key)
			}
		}
		Mutex.Unlock()
		return fmt.Errorf("The key [%s] has reached max requests [%d]", key, k.Hits)
	}
	k.Hit()
	Mutex.Unlock()
	return nil
}

func (r *RateLimit) Hit() {
	r.Hits += 1
}
