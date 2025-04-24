package mqtt

import (
    "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

func Init() error {
    opts := mqtt.NewClientOptions()

    opts.AddBroker("tcp://localhost:1883")
    opts.SetClientID("pomodoro-server")

    Client = mqtt.NewClient(opts)

    if token := Client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }

    return nil
}
