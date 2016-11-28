package timer

import (
	"container/heap"
	"time"
)

type _Timer struct {
	fireTime time.Time
	interval time.Duration
	callback CallbackFunc
	repeat   bool
}

type _TimerHeap struct {
	timers []_Timer
}

func (h *_TimerHeap) Len() int {
	return len(h.timers)
}

func (h *_TimerHeap) Less(i, j int) bool {
	return h.timers[i].fireTime.Before(h.timers[j].fireTime)
}

func (h *_TimerHeap) Swap(i, j int) {
	var tmp _Timer
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *_TimerHeap) Push(x interface{}) {
	h.timers = append(h.timers, x.(_Timer))
}

func (h *_TimerHeap) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

// Type of callback function
type CallbackFunc func()

var (
	timerHeap _TimerHeap
)

func init() {
	heap.Init(&timerHeap)
}

// Add a callback which will be called after specified duration
func AddCallback(t time.Duration, callback CallbackFunc) {
	heap.Push(&timerHeap, _Timer{
		fireTime: time.Now().Add(t),
		interval: t,
		callback: callback,
		repeat:   false,
	})
}

// Add a timer which calls callback periodly
func AddTimer(t time.Duration, callback CallbackFunc) {
	heap.Push(&timerHeap, _Timer{
		fireTime: time.Now().Add(t),
		interval: t,
		callback: callback,
		repeat:   true,
	})
}

// Tick once for timers
func Tick() {
	now := time.Now()
	for {
		if timerHeap.Len() <= 0 {
			break
		}

		nextFireTime := timerHeap.timers[0].fireTime
		if nextFireTime.After(now) {
			break
		}
		t := heap.Pop(&timerHeap).(_Timer)
		t.callback()
		if t.repeat {
			// add Timer back to heap
			t.fireTime = t.fireTime.Add(t.interval)
			if !t.fireTime.After(now) {
				t.fireTime = now.Add(t.interval)
			}
			heap.Push(&timerHeap, t)
		}
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
