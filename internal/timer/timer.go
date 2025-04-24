package timer

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/lukepetko/pomodoro-server/internal/mqtt"
)

type Timer struct {
    duration int
    remaining int
    startedAt time.Time
    running bool
    lock sync.Mutex
    stopChan chan bool
    doneChan chan bool
}

func New(duration int) *Timer {
    return &Timer{
        duration: duration,
        remaining: duration,
        startedAt: time.Now(),
        running: false,
        stopChan: make(chan bool),
        doneChan: make(chan bool),
    }
}

func (t *Timer) Start() {
    t.lock.Lock()
    if t.running {
        t.lock.Unlock()
        return
    }
    t.running = true
    t.lock.Unlock()

    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            t.lock.Lock()
            mqtt.Client.Publish("pomodoro/timer/tick", 0, false, strconv.Itoa(t.remaining))
            fmt.Println(t.remaining)
            t.remaining--
            if t.remaining <= 0 {
                t.lock.Unlock()
                t.doneChan <- true
                return
            }
            t.lock.Unlock()
        }
    } ()
}

func (t *Timer) Done() <-chan bool {
    return t.doneChan
}

