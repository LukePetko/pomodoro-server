package api

import (
    "net/http"
	"github.com/lukepetko/pomodoro-server/internal/timer"
)

type Server struct {
    timer *timer.Timer
}

func NewServer(timer *timer.Timer) *Server {
    return &Server{
        timer: timer,
    }
}

func (s *Server) StartTimer(w http.ResponseWriter, r *http.Request) {
    s.timer.Start()
    w.Write([]byte("Timer started"))
}

func (s *Server) StopTimer(w http.ResponseWriter, r *http.Request) {
    s.timer.Stop()
    w.Write([]byte("Timer stopped"))
}

func (s *Server) RestartTimer(w http.ResponseWriter, r *http.Request) {
    s.timer.Restart()
    w.Write([]byte("Timer restarted"))
}

func (s *Server) Routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/start", s.StartTimer)
    mux.HandleFunc("/stop", s.StopTimer)
    mux.HandleFunc("/restart", s.RestartTimer)

    return mux
}
