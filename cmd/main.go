package main

import (
	"fmt"

	config "github.com/lukepetko/pomodoro-server/internal/config"
	mqtt "github.com/lukepetko/pomodoro-server/internal/mqtt"
	"github.com/lukepetko/pomodoro-server/internal/timer"
)



func main() {
    fmt.Println("MQTT Time Tracker app started!")
    
    if err := mqtt.Init(); err != nil {
        fmt.Println(err)
        return
    }

    config, err := config.LoadConfig("config.json")

    if err != nil {
        fmt.Println(err)
        return
    }

    timer := timer.New(config.WorkTime)

    timer.Start()

    <-timer.Done()
}

