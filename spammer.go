package ratelimit

import "time"

type Spam struct {
	ExpiredAt time.Time
	Hits      int
}

type Spammer struct {
	Duration time.Duration
	Values   map[string]*Spam
}

func CreateSpammer() Spammer {
	sp := Spammer{
		Duration: time.Hour * 24,
		Values:   make(map[string]*Spam),
	}
	SpamCleaner(&sp)
	return sp
}

func (s Spammer) Increase(key string) {
	k, ok := s.Values[key]
	if !ok {
		s.Values[key] = createSpam(s.Duration)
		k = s.Values[key]
	}
	k.Hits += 1
}

func createSpam(d time.Duration) *Spam {
	return &Spam{
		ExpiredAt: time.Now().Add(d),
		Hits:      0,
	}
}

func SpamCleaner(l *Spammer) {
	go func(l *Spammer) {
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
