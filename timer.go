package timer

import (
	"container/heap"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

const (
	MIN_TIMER_INTERVAL = 1 * time.Millisecond
)

type Timer struct {
	fireTime  time.Time
	interval  time.Duration
	callback  CallbackFunc
	repeat    bool
	cancelled bool
}

func (t *Timer) Cancel() {
	t.cancelled = true
}

func (t *Timer) IsActive() bool {
	return !t.cancelled
}

type _TimerHeap struct {
	timers []*Timer
}

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

func (h *_TimerHeap) Less(i, j int) bool {
	return h.timers[i].fireTime.Before(h.timers[j].fireTime)
}

func (h *_TimerHeap) Swap(i, j int) {
	var tmp *Timer
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *_TimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(*Timer))
}

func (h *_TimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

// Type of callback function
type CallbackFunc func()

var (
	timerHeap     _TimerHeap
	timerHeapLock sync.RWMutex
)

func init() {
	heap.Init(&timerHeap)
}

// Add a callback which will be called after specified duration
func AddCallback(d time.Duration, callback CallbackFunc) *Timer {
	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		repeat:   false,
	}
	timerHeapLock.Lock()
	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

// Add a timer which calls callback periodly
func AddTimer(d time.Duration, callback CallbackFunc) *Timer {
	if d < MIN_TIMER_INTERVAL {
		d = MIN_TIMER_INTERVAL
	}

	t := &Timer{
		fireTime: time.Now().Add(d),
		interval: d,
		callback: callback,
		repeat:   true,
	}
	timerHeapLock.Lock()
	heap.Push(&timerHeap, t)
	timerHeapLock.Unlock()
	return t
}

// Tick once for timers
func Tick() {
	now := time.Now()
	isWriteLock := false
	timerHeapLock.RLock()

	for {
		if timerHeap.Len() <= 0 {
			break
		}

		nextFireTime := timerHeap.timers[0].fireTime
		//fmt.Printf(">>> nextFireTime %s, now is %s\n", nextFireTime, now)
		if nextFireTime.After(now) {
			break
		}
		// require a write lock since then
		if !isWriteLock {
			timerHeapLock.RUnlock()
			timerHeapLock.Lock()
			isWriteLock = true
		}
		t := heap.Pop(&timerHeap).(*Timer)

		if t.cancelled {
			continue
		}

		if !t.repeat {
			t.cancelled = true
		}

		runCallback(t.callback)

		if t.repeat {
			// add Timer back to heap
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) {
				t.fireTime = now.Add(t.interval)
			}
			heap.Push(&timerHeap, t)
		}
	}
	if !isWriteLock {
		timerHeapLock.RUnlock()
	} else {
		timerHeapLock.Unlock()
	}
}

// Start the self-ticking routine, which ticks per tickInterval
func StartTicks(tickInterval time.Duration) {
	go selfTickRoutine(tickInterval)
}

func selfTickRoutine(tickInterval time.Duration) {
	for {
		time.Sleep(tickInterval)
		Tick()
	}
}

func runCallback(callback CallbackFunc) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback %v paniced: %v\n", callback, err)
			debug.PrintStack()
		}
	}()
	callback()
}
