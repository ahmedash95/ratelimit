package ratelimit

import "time"

type Block struct {
	ExpiredAt time.Time
}

type Blocker struct {
	Duration time.Duration
	Values   map[string]*Block
}

func CreateBlocker() Blocker {
	b := Blocker{
		Duration: time.Hour * 24,
		Values:   make(map[string]*Block),
	}
	BlockerCleaner(&b)
	return b
}

func (s Blocker) AddIfNotExists(key string) {
	_, ok := s.Values[key]
	if !ok {
		s.Values[key] = createBlock(s.Duration)
	}
}

func createBlock(d time.Duration) *Block {
	return &Block{
		ExpiredAt: time.Now().Add(d),
	}
}

func BlockerCleaner(l *Blocker) {
	go func(l *Blocker) {
		for {
			time.Sleep(time.Second)
			now := time.Now()
			Mutex.Lock()
			for k, r := range l.Values {
				if r.ExpiredAt.Before(now) {
					delete(l.Values, k)
				}
			}
			Mutex.Unlock()
		}
	}(l)
}
