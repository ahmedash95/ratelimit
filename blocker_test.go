package ratelimit

import (
	"testing"
	"time"

	"aahframework.org/test.v0/assert"
	"bou.ke/monkey"
)

func Test_it_block_as_expected(t *testing.T) {
	l1 := CreateLimit("1r/s,spam:3,block:2d")
	key := "127.0.0.1"
	l1.Hit(key)
	l1.Hit(key)
	time.Sleep(1 * time.Second)
	l1.Hit(key)
	l1.Hit(key)
	time.Sleep(1 * time.Second)
	l1.Hit(key)
	l1.Hit(key)
	assert.NotNil(t, l1.Blocker.Values[key])
}

func Test_it_clears_block_after_expeced_duration(t *testing.T) {
	l1 := CreateLimit("1r/s,spam:3,block:2d")
	key := "127.0.0.1"
	l1.Hit(key)
	l1.Hit(key)

	travelTowDayWithOneHour := time.Now().Add(49 * time.Hour)
	patch := monkey.Patch(time.Now, func() time.Time { return travelTowDayWithOneHour })
	defer patch.Unpatch()

	time.Sleep(2 * time.Second)

	assert.Nil(t, l1.Blocker.Values[key])
}
