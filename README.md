# RateLimit
A high-performance rate limiter written in GO language

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![godoc](https://godoc.org/github.com/ahmedash95/ratelimit?status.svg)](https://godoc.org/github.com/ahmedash95/ratelimit)
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmedash95/ratelimit)](https://goreportcard.com/report/github.com/ahmedash95/ratelimit)

This package provides a Golang implementation of rate limit

Create a rate limiter with a maximum number of operations to perform per second. with the ability of mark as a spam then block for a specific period

## Usage

```golang
package main

import (
	"fmt"

	"github.com/ahmedash95/ratelimit"
)

func main() {
    rt := ratelimit.CreateLimit("1r/s")
    user_ip := "127.0.0.1"
    rt.Hit(user_ip)

    fmt.Println(rt.Rates[user_ip].Hits)
    // out: 1
}
```

So now how you would block the user if he hits the request more than the limit

```golang
package main

import (
	"fmt"

	"github.com/ahmedash95/ratelimit"
)

func main() {
	rt := ratelimit.CreateLimit("1r/s")
	user_ip := "127.0.0.1"
	var err error
	err = rt.Hit(user_ip)
	err = rt.Hit(user_ip)

    if err != nil {
        fmt.Println(err)
        // out: The key [127.0.0.1] has reached max requests [1]

        // block process
    }
}
```

### Define spam and block method

the way of how it works is `1r/s,spam:3,block:3d` means that the rate limit is 1 request per second and we will mark it as a spammer if the `key` reach the max limit more than `3 times`, if he, we will block the `key` for `3 days`

```golang
package main

import (
	"fmt"
	"github.com/ahmedash95/ratelimit"
)

func main() {
	rt := ratelimit.CreateLimit("1r/s,spam:3,block:2d")
	user_ip := "127.0.0.1"
	rt.Hit(user_ip)
    rt.Hit(user_ip)

    fmt.Println(rt.Spammer.Values[user_ip].Hits)
    // out: 1
    // because he just hit the max requests for 1 time

    rt := ratelimit.CreateLimit("1r/s,spam:3,block:2d")
    user_ip := "127.0.0.1"
    rt.Hit(user_ip)
    rt.Hit(user_ip)

    time.Sleep(time.Second)
    rt.Hit(user_ip)
    rt.Hit(user_ip)

    time.Sleep(time.Second)

    rt.Hit(user_ip)
    rt.Hit(user_ip)

    fmt.Println(rt.Spammer.Values[user_ip].Hits)
    // out: 3
    // because he hit the max requests for 3 times

    _, blocked := rt.Blocker.Values[user_ip]
    fmt.Printf("%v\n", blocked)
    // out: True
}
```

## Guide

### Rate limit pattern definition

the default format is `1r/s` which is mean 1 request per second or `3r/m` which is mean 3 request per minute

**The advanced Pattern that support spam and block**

`1r/s,spam:3,block:3h` which is mean block the user for 3 hours when reatchs the maximum spam hits, **block** supports `3d` `10h` `5m` `10s`

## Testing

```
$ go test
```

## Credits

- [Ahmed Ashraf](https://github.com/ahmedash95)

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
