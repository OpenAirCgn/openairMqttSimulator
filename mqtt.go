package openairMqttSimulator

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// creates a new (acutal, not simulated) mqtt client, connecting to the
// indicated URL using the provided client id.
func newMqttClient(url string, clientId string) (mqtt.Client, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(url)
	options.SetClientID(clientId)
	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

type MqttClientSim struct {
	BrokerUrl   string        // tcp://bla:1883
	NumClients  int           // number of client to run simultaneously
	Frequency   time.Duration // how long to sleep between "measurements"
	NumRequests int           // total number of requests to send
	LogRequests bool
}

func (cl *MqttClientSim) RunSimulation() {
	var wg sync.WaitGroup
	for i := 0; i != cl.NumClients; i++ {
		wg.Add(1)
		go cl.runSimulation(&wg, i)
	}
	wg.Wait()

}

func sleepJitter(d time.Duration, amt float32) {
	var sign float32
	if rand.NormFloat64() > 0 {
		sign = 1.0
	} else {
		sign = -1.0
	}
	jitter := (time.Duration)(sign * float32(d) * amt * rand.Float32())

	time.Sleep(d + jitter)

}

func (cl *MqttClientSim) runSimulation(wg *sync.WaitGroup, clientNum int) {
	defer wg.Done()

	device := NewQuadsenseDevice(int32(clientNum))
	client, err := newMqttClient(cl.BrokerUrl, device.DeviceId.String())
	if err != nil {
		panic(err)
	}
	for i := 0; i != cl.NumRequests; i++ {
		for _, sensor := range device.Sensors {
			value := sensor.Value()
			if cl.LogRequests {
				fmt.Fprintf(os.Stderr, "%s t: %s msg: %s\n", time.Now().Format("2006-01-02T15:04:05.999"), device.Topic(sensor), value)
			}
			client.Publish(device.Topic(sensor), 0, false, value)
			sleepJitter(10*time.Millisecond, 0.1)

		}
		sleepJitter(cl.Frequency, 0.01)
	}
}
