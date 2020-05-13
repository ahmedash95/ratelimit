package ratelimit

import (
	"testing"
	"time"

	"aahframework.org/test.v0/assert"
	"bou.ke/monkey"
)

func Test_it_count_spam_as_expected(t *testing.T) {
	l1 := CreateLimit("1r/s,spam:3,block:2d")
	key := "127.0.0.1"
	l1.Hit(key)
	l1.Hit(key)
	expected := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
	time.Sleep(1*time.Second + (time.Millisecond * 100))
	l1.Hit(key)
	l1.Hit(key)
	time.Sleep(1*time.Second + (time.Millisecond * 100))
	l1.Hit(key)
	l1.Hit(key)
	assert.Equal(t, 3, l1.Spammer.Values[key].Hits)
	actual := l1.Spammer.Values[key].ExpiredAt.Format("2006-01-02 15:04:05")
	assert.Equal(t, expected, actual)
}

func Test_it_clears_spam_after_expeced_duration(t *testing.T) {
	l1 := CreateLimit("1r/s,spam:3,block:2d")
	key := "127.0.0.1"
	l1.Hit(key)
	l1.Hit(key)

	travelOneDayWithOneHour := time.Now().Add(25 * time.Hour)
	patch := monkey.Patch(time.Now, func() time.Time { return travelOneDayWithOneHour })
	defer patch.Unpatch()

	time.Sleep(2 * time.Second)

	assert.Nil(t, l1.Spammer.Values[key])
}
