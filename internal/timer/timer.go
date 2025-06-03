package timer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	config "github.com/lukepetko/pomodoro-server/internal/config"
	mqtt "github.com/lukepetko/pomodoro-server/internal/mqtt"
)

type Timer struct {
	duration  int
	remaining int
    session   int
    sessions  []int
	running   bool
	lock      sync.Mutex
	stopChan  chan bool
	doneChan  chan bool
}

type SessionMessage struct {
    SessionNumber int    `json:"session_number"`
    TimerType     string `json:"timer_type"`
    EventType     string `json:"event_type"`
    Running       bool   `json:"running"`
    Duration      int    `json:"duration"`
}

type TimerState struct {
	Duration  int
	Remaining int
	Session   int
	Sessions  []int
    Running   bool
}

func GetTimerType(session int, sessions []int) string {
    if session % 2 == 0 {
        return "work"
    } else if session == len(sessions) - 1 {
        return "long_break"
    } else {
        return "short_break"
    }
}

func New(config *config.Config) *Timer {
    var sessions []int
    for i := 0; i < config.NumberOfSessions; i++ {
        sessions = append(sessions, config.WorkTime)
        if i == config.NumberOfSessions - 1 {
            sessions = append(sessions, config.LongBreakTime)
        } else {
            sessions = append(sessions, config.ShortBreakTime)
        }
    }
	return &Timer{
		duration:  sessions[0],
		remaining: sessions[0],
        session:   0,
        sessions:  sessions,
		running:   false,
		stopChan:  make(chan bool),
		doneChan:  make(chan bool),
	}
}

func (t *Timer) StartProcess() {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            t.lock.Lock()
            if t.running {
                mqtt.Client.Publish("pomodoro/timer/tick", 0, false, strconv.Itoa(t.remaining))
                fmt.Println(t.remaining)
                t.remaining--
            }
            if t.remaining <= 0 && t.session < len(t.sessions) - 1 {
                fmt.Println("Session complete " + strconv.Itoa(t.session) + " remaining " + strconv.Itoa(t.remaining))
                payload := SessionMessage{
                    SessionNumber: t.session,
                    TimerType:     GetTimerType(t.session, t.sessions),
                    EventType:     "end",
                    Running:       t.running,
                    Duration:      t.duration,
                }
                jsonPayload, _ := json.Marshal(payload)
                mqtt.Client.Publish("pomodoro/timer/session", 0, false, string(jsonPayload))
                t.session++
                t.remaining = t.sessions[t.session]
                t.duration = t.remaining
                payload = SessionMessage{
                    SessionNumber: t.session,
                    TimerType:     GetTimerType(t.session, t.sessions),
                    EventType:     "start",
                    Running:       t.running,
                    Duration:      t.duration,
                }
                jsonPayload, _ = json.Marshal(payload)
                mqtt.Client.Publish("pomodoro/timer/session", 0, false, string(jsonPayload))
            } else if t.remaining <= 0 && t.session == len(t.sessions) - 1 {
                payload := SessionMessage{
                    SessionNumber: t.session,
                    TimerType:     "long_break",
                    EventType:     "end",
                    Running:       t.running,
                    Duration:      t.duration,
                }
                jsonPayload, _ := json.Marshal(payload)
                mqtt.Client.Publish("pomodoro/timer/session", 0, false, string(jsonPayload))
                t.session = 0
                t.remaining = t.sessions[t.session]
                t.duration = t.remaining
                t.running = false
            }
            t.lock.Unlock()
        }
    }()
}

func (t *Timer) Start() {
	t.lock.Lock()
	if t.running {
		t.lock.Unlock()
		return
	}
	t.running = true
	t.lock.Unlock()

}

func (t *Timer) Stop() {
	t.lock.Lock()
	if !t.running {
		t.lock.Unlock()
		return
	}
	t.running = false
	t.lock.Unlock()

	go func() {
		t.stopChan <- true
	}()
}

func (t *Timer) Restart() {
    t.lock.Lock()
    t.session = 0
    t.remaining = t.sessions[t.session]
    t.lock.Unlock()
}

func (t *Timer) Status() TimerState {
    t.lock.Lock()
	defer t.lock.Unlock()

	return TimerState{
		Duration:  t.duration,
		Remaining: t.remaining,
		Session:   t.session,
		Sessions:  t.sessions,
        Running:   t.running,
	}
}

func (t *Timer) Done() <-chan bool {
	return t.doneChan
}
