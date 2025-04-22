package main

import (
	"fmt"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
    fmt.Println("MQTT Time Tracker app started!")
    
    opts := mqtt.NewClientOptions()
    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("pomodoro-server")

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
    
    number := 0

    for {
        time.Sleep(time.Second)
        fmt.Println("Sending message", number)
        client.Publish("test/test", 0, false, strconv.Itoa(number))
        number++
    }
}

