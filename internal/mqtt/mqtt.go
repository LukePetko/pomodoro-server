package mqtt

import (
    "os"
    "github.com/eclipse/paho.mqtt.golang"
)

var Client mqtt.Client

func Init() error {
    opts := mqtt.NewClientOptions()

    address := os.Getenv("MQTT_BROKER_ADDRESS")
    port := os.Getenv("MQTT_BROKER_PORT")

    opts.AddBroker("tcp://" + address + ":" + port)
    opts.SetClientID("pomodoro-server")

    Client = mqtt.NewClient(opts)

    if token := Client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }

    return nil
}
