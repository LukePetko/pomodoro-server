package main

import (
	"fmt"
    "net/http"
    "log"

	config "github.com/lukepetko/pomodoro-server/internal/config"
	mqtt "github.com/lukepetko/pomodoro-server/internal/mqtt"
	"github.com/lukepetko/pomodoro-server/internal/timer"
    "github.com/lukepetko/pomodoro-server/internal/api"
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

    timer := timer.New(config)
    timer.StartProcess()

    srv := api.NewServer(timer, config)

    fmt.Println("Server started at port 9200!")
    log.Fatal(http.ListenAndServe(":9200", srv.Routes()))
}

