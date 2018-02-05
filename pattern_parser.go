package ratelimit

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Max       int
	Per       time.Duration
	Block     time.Duration
	MaxToSpam int
}

func parse(s string) (Options, error) {
	op := Options{}

	parsed := strings.Split(s, ",")

	p := strings.Split(parsed[0], "r/")

	maxValue, _ := strconv.Atoi(p[0])
	if maxValue < 1 {
		return op, fmt.Errorf("Invalid max requests value [%d] must be larger than zero", maxValue)
	}

	per := p[1]
	if per != "s" && per != "m" {
		return op, fmt.Errorf("Invalid limit [%s] Limit must be per s as second or m as minute", per)
	}

	perDuration := time.Second
	if per == "m" {
		perDuration = time.Minute
	}

	op.Max = maxValue
	op.Per = perDuration

	for _, p := range parsed[1:] {
		s := strings.Split(p, ":")
		if len(s) != 2 {
			return op, fmt.Errorf("Can't parse value: %s", p)
		}
		if s[0] != "spam" && s[0] != "block" {
			return op, fmt.Errorf("Unsupported module [%s] must be spam or block", s[0])
		}
		if s[0] == "spam" {
			value, _ := strconv.Atoi(s[1])
			op.MaxToSpam = value
		}
		if s[0] == "block" {
			v := s[1]
			max, _ := strconv.Atoi(v[:len(v)-1])
			duration := v[len(v)-1:]
			if duration != "d" && duration != "h" && duration != "m" && duration != "s" {
				return op, fmt.Errorf("Unsupported time duration [%s] must be (d) for day, (h) for hour, (m) for minute or (s) for second.", duration)
			}
			op.Block = getDuration(max, duration)
		} else {
			op.Block = 0
		}
	}
	return op, nil
}

func getDuration(max int, s string) time.Duration {
	t := time.Second
	switch s {
	case "d":
		t = time.Hour * 24
	case "h":
		t = time.Hour
	case "m":
		t = time.Minute
	default:
		t = time.Second
	}
	return time.Duration(max) * t
}
