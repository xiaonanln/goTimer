package timer_test

import (
	"fmt"
	"testing"
	"time"

	"math/rand"

	"log"

	"runtime/pprof"

	"os"

	"github.com/xiaonanln/goTimer"
)

func init() {
	timer.StartTicks(time.Millisecond)
}

func TestCallback(t *testing.T) {
	INTERVAL := 100 * time.Millisecond
	for i := 0; i < 10; i++ {
		x := false
		timer.AddCallback(INTERVAL, func() {
			fmt.Println("callback!")
			x = true
		})
		time.Sleep(INTERVAL * 2)
		if !x {
			t.Fatalf("x should be true, but it's false")
		}
	}
}

func TestTimer(t *testing.T) {
	INTERVAL := 100 * time.Millisecond
	x := 0
	px := x
	now := time.Now()
	nextTime := now.Add(INTERVAL)
	fmt.Printf("now is %s, next time should be %s\n", time.Now(), nextTime)

	timer.AddTimer(INTERVAL, func() {
		x += 1
		fmt.Printf("timer %s x %v px %v\n", time.Now(), x, px)
	})
	//time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		time.Sleep(nextTime.Add(INTERVAL / 2).Sub(time.Now()))
		fmt.Printf("Check x %v px %v @ %s\n", x, px, time.Now())
		if x != px+1 {
			t.Fatalf("x should be %d, but it's %d", px+1, x)
		}
		px = x
		nextTime = nextTime.Add(INTERVAL)
		fmt.Printf("now is %s, next time should be %s\n", time.Now(), nextTime)
	}
}

func TestCallbackSeq(t *testing.T) {
	a := 0
	d := time.Second

	for i := 0; i < 100; i++ {
		i := i
		timer.AddCallback(d, func() {
			if a != i {
				t.Error(i, a)
			}

			a += 1
		})
	}
	time.Sleep(d + time.Second*1)
}

func TestCancelCallback(t *testing.T) {
	INTERVAL := 20 * time.Millisecond
	x := 0

	timer := timer.AddCallback(INTERVAL, func() {
		x = 1
	})
	if !timer.IsActive() {
		t.Fatalf("timer should be active")
	}
	timer.Cancel()
	if timer.IsActive() {
		t.Fatalf("timer should be inactive")
	}
	time.Sleep(INTERVAL * 2)
	if x != 0 {
		t.Fatalf("x should be 0, but is %v", x)
	}
}

func TestCancelTimer(t *testing.T) {
	INTERVAL := 20 * time.Millisecond
	x := 0
	timer := timer.AddTimer(INTERVAL, func() {
		x += 1
	})
	if !timer.IsActive() {
		t.Fatalf("timer should be active")
	}
	timer.Cancel()
	if timer.IsActive() {
		t.Fatalf("timer should be inactive")
	}
	time.Sleep(INTERVAL * 2)
	if x != 0 {
		t.Fatalf("x should be 0, but is %v", x)
	}
}

func TestTimerPerformance(t *testing.T) {
	f, err := os.Create("TestTimerPerformance.cpuprof")
	if err != nil {
		panic(err)
	}

	pprof.StartCPUProfile(f)
	duration := 10 * time.Second

	for i := 0; i < 400000; i++ {
		if rand.Float32() < 0.5 {
			d := time.Duration(rand.Int63n(int64(duration)))
			timer.AddCallback(d, func() {})
		} else {
			d := time.Duration(rand.Int63n(int64(time.Second)))
			timer.AddTimer(d, func() {})
		}
	}

	log.Println("Waiting for", duration, "...")
	time.Sleep(duration)
	pprof.StopCPUProfile()
}
