package util

//https://pjchender.dev/golang/pkg-time-rate/#waitwaitn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestWait(*testing.T) {
	counter := 0
	ctx := context.Background()

	limit := rate.Every(time.Millisecond * 200)
	limiter := rate.NewLimiter(limit, 1)
	fmt.Println(limiter.Limit(), limiter.Burst()) // 5，1

	for {
		counter++
		limiter.Wait(ctx)
		fmt.Printf("counter: %v, %v \n", counter, time.Now().Format(time.RFC3339))
	}
}

func TestAllow(*testing.T) {
	counter := 0

	limit := rate.Every(time.Millisecond * 200)
	limiter := rate.NewLimiter(limit, 4)
	fmt.Println(limiter.Limit(), limiter.Burst()) // 5，4

	for {
		counter++

		if isAllowed := limiter.AllowN(time.Now(), 3); isAllowed {
			fmt.Printf("counter: %v, %v \n", counter, time.Now().Format(time.RFC3339))
		} else {
			fmt.Printf("counter: %v, not allow \n", counter)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func TestReserve(*testing.T) {
	counter := 0

	limit := rate.Every(time.Millisecond * 200)
	limiter := rate.NewLimiter(limit, 3)
	fmt.Println(limiter.Limit(), limiter.Burst()) // 5，3

	for {
		counter++
		tokensNeed := 2
		reserve := limiter.ReserveN(time.Now(), tokensNeed)

		if !reserve.OK() {
			fmt.Printf("required token(%v) once is greater thab capacity token count(%v)\n", tokensNeed, limiter.Burst())
			return
		}

		time.Sleep(reserve.Delay())

		fmt.Printf("counter: %v, %v \n", counter, time.Now().Format(time.RFC3339))
	}
}
