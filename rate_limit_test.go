package ratelimit

import (
	"fmt"
	"testing"
	"time"

	"aahframework.org/test.v0/assert"
	"bou.ke/monkey"
)

func Test_IT_PANIC_WHEN_INVALID_NEW_RATE_PATTERN(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The CreateLimit did not panic when invalid pattern used")
		}
	}()

	_ = CreateLimit("1t/s")
}

func Test_IT_PANIC_WHEN_INVALID_TIME_PATTERN(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The CreateLimit did not panic when invalid time pattern used")
		}
	}()

	_ = CreateLimit("1r/d")
}

func Test_HIT_New_Rate_Limits(t *testing.T) {
	l1 := CreateLimit("1r/s")
	k1 := "127.0.0.1"
	l1.Hit(k1)
	assert.Equal(t, 1, l1.Rates[k1].Hits)
}

func Test_Rate_Limits_WITH_OPTIONS(t *testing.T) {
	l1 := CreateLimit("1r/s,spam:3,block:2d")
	assert.Equal(t, time.Second, l1.Per)
	assert.Equal(t, 1, l1.MaxRequests)
	assert.Equal(t, 2*time.Hour*24, l1.Block)
}

func Test_IT_RETURN_ERROR_AFTER_HIT_OVER_LIMIT(t *testing.T) {
	l1 := CreateLimit("1r/s")
	k1 := "127.0.0.1"
	var err error
	for {
		if err != nil {
			break
		}
		err = l1.Hit(k1)
	}

	expected := fmt.Sprintf("The key [%s] has reached max requests [1]", k1)
	actual := err

	if actual.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}

func Test_It_Clear_The_Limit_EVERY_SECOND(t *testing.T) {
	l := CreateLimit("1r/s")
	k := "127.0.0.1"

	l.Hit(k)
	time.Sleep(2 * time.Second)

	l.Hit(k)
	time.Sleep(2 * time.Second)

	assert.Equal(t, 0, l.Rates[k].Hits)
}

func Test_It_Clear_The_Limit_EVERY_MINUTE(t *testing.T) {
	l := CreateLimit("5r/m")
	k := "127.0.0.1"
	l.Hit(k)
	l.Hit(k)
	l.Hit(k)

	assert.Equal(t, 3, l.Rates[k].Hits)

	travelOneMinute := time.Now().Add(time.Minute)
	patch := monkey.Patch(time.Now, func() time.Time { return travelOneMinute })
	defer patch.Unpatch()

	time.Sleep(2 * time.Second)
	l.Hit(k)
	assert.Equal(t, 1, l.Rates[k].Hits)
}
