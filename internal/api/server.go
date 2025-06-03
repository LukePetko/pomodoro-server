package api

import (
	"encoding/json"
	"fmt"
	"github.com/lukepetko/pomodoro-server/internal/config"
	"github.com/lukepetko/pomodoro-server/internal/timer"
	"net/http"

	mqtt "github.com/lukepetko/pomodoro-server/internal/mqtt"
)

type Server struct {
	timer  *timer.Timer
	config *config.Config
}

func NewServer(timer *timer.Timer, config *config.Config) *Server {
	return &Server{
		timer:  timer,
		config: config,
	}
}

func (s *Server) StartTimer(w http.ResponseWriter, r *http.Request) {
	s.timer.Start()
    status := s.timer.Status()
	payload := timer.SessionMessage{
        SessionNumber: status.Session,
        TimerType:     timer.GetTimerType(status.Session, status.Sessions),
		EventType:     "start",
        Running:       status.Running,
        Duration:      status.Duration,
	}
	jsonPayload, _ := json.Marshal(payload)
	mqtt.Client.Publish("pomodoro/timer/session", 0, false, string(jsonPayload))
	w.Write([]byte("Timer started"))
}

func (s *Server) StopTimer(w http.ResponseWriter, r *http.Request) {
	s.timer.Stop()
    status := s.timer.Status()
	payload := timer.SessionMessage{
		SessionNumber: status.Session,
		TimerType:     timer.GetTimerType(status.Session, status.Sessions),
		EventType:     "stop",
		Running:       status.Running,
        Duration:      status.Duration,
	}
	jsonPayload, _ := json.Marshal(payload)
	mqtt.Client.Publish("pomodoro/timer/session", 0, false, string(jsonPayload))
	w.Write([]byte("Timer stopped"))
}

func (s *Server) RestartTimer(w http.ResponseWriter, r *http.Request) {
	s.timer.Restart()
    status := s.timer.Status()
	payload := timer.SessionMessage{
        SessionNumber: status.Session,
        TimerType:     timer.GetTimerType(status.Session, status.Sessions),
		EventType:     "restart",
        Running:       status.Running,
        Duration:      status.Duration,
	}
	jsonPayload, _ := json.Marshal(payload)
	mqtt.Client.Publish("pomodoro/timer/session", 0, false, string(jsonPayload))
	w.Write([]byte("Timer restarted"))
}

func (s *Server) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	state := s.timer.Status()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (s *Server) SaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Saving config")

	var newCfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Println(newCfg)

	if err := config.SaveConfig("config.json", &newCfg); err != nil {
		http.Error(w, "Failed to save config", http.StatusInternalServerError)
		return
	}

	*s.config = newCfg
	s.timer.Restart()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Config updated and saved"})
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/start", s.StartTimer)
	mux.HandleFunc("/stop", s.StopTimer)
	mux.HandleFunc("/restart", s.RestartTimer)
	mux.HandleFunc("/config", s.SaveConfig)
	mux.HandleFunc("/status", s.Status)

	return mux
}
