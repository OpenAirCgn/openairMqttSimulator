package openairMqttSimulator

import (
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

//func TestMqtt(t *testing.T) {
//	sim := MqttClientSim{
//		10, 1 * time.Second, 10, nil,
//	}
//	sim.RunSimulation()
//}

func TestMqttForRealz(t *testing.T) {
	options := mqtt.NewClientOptions()
	options.AddBroker("tcp://test.mosquitto.org:1883")
	options.SetClientID("00:00:00:00:00:00")
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(100)
	sim := MqttClientSim{
		5,
		100 * time.Millisecond,
		10,
		client,
	}

	sim.RunSimulation()
}
