package main

import (
	"fmt"

	config "github.com/lukepetko/pomodoro-server/internal/config"
	mqtt "github.com/lukepetko/pomodoro-server/internal/mqtt"
	"github.com/lukepetko/pomodoro-server/internal/timer"
    "github.com/joho/godotenv"
)



func main() {
    fmt.Println("MQTT Time Tracker app started!")

    err := godotenv.Load()
    if err != nil {
        fmt.Println(err)
        return
    }
    
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

