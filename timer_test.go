package timer

import (
	"fmt"
	"testing"
	"time"
)

func init() {
	StartTicks(time.Millisecond)
}

func TestCallback(t *testing.T) {
	for i := 0; i < 10; i++ {
		x := false
		AddCallback(10*time.Millisecond, func() {
			fmt.Println("callback!")
			x = true
		})
		time.Sleep(20 * time.Millisecond)
		if !x {
			t.Fatalf("x should be true, but it's false")
		}
	}
}

func TestTimer(t *testing.T) {
	x := 0
	px := x
	nextTime := time.Now().Add(10 * time.Millisecond)
	fmt.Printf("now is %s, next time should be %s\n", time.Now(), nextTime)

	AddTimer(10*time.Millisecond, func() {
		x += 1
		fmt.Printf("timer %s x %v px %v\n", time.Now(), x, px)
	})

	for i := 0; i < 100; i++ {
		time.Sleep(nextTime.Add(5 * time.Millisecond).Sub(time.Now()))
		fmt.Printf("Check x %v px %v @ %s", x, px, time.Now())
		if x != px+1 {
			t.Fatalf("x should be %d, but it's %d", px+1, x)
		}
		px = x
		nextTime = nextTime.Add(10 * time.Millisecond)
		fmt.Printf("now is %s, next time should be %s\n", time.Now(), nextTime)
	}
}
