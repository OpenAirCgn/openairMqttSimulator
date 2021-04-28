package openairMqttSimulator

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClientSim struct {
	BrokerUrl     string        // tcp://bla:1883 tls://... ws:// wss://
	TLSConfig     *tls.Config   // optional TLS configuration
	ClientCertDir string        // directory to search for tls client certificates.
	NumClients    int           // number of client to run simultaneously
	Frequency     time.Duration // how long to sleep between "measurements"
	NumRequests   int           // total number of requests to send
	LogRequests   bool          // whether or not to be verbose
	UseSha        bool          // use sha1 hash of mac instead of mac
	UseCounter    bool          // use a counter as the first value transmitted
	Qos           byte
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

func setClientCertIfAvailable(clientId string, dir string, cfg *tls.Config) {
	crtFN := fmt.Sprintf("%s/%s.crt", dir, clientId)
	keyFN := fmt.Sprintf("%s/%s.pem", dir, clientId)

	cert, err := tls.LoadX509KeyPair(crtFN, keyFN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load client cert for %s (%v)\n", clientId, err)
		return
	}
	cfg.Certificates = []tls.Certificate{cert}
}

// creates a new (acutal, not simulated) mqtt client, connecting to the
// indicated URL using the provided client id.
func (cl *MqttClientSim) newMqttClient(clientId string) (mqtt.Client, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(cl.BrokerUrl)
	options.SetClientID(clientId)

	if strings.HasPrefix(cl.BrokerUrl, "tls://") || strings.HasPrefix(cl.BrokerUrl, "wss://") {
		cfg := cl.TLSConfig
		if cfg == nil {
			cfg = &tls.Config{}
		}

		// attempt to set a client certificate if available
		setClientCertIfAvailable(clientId, cl.ClientCertDir, cfg)

		options.SetTLSConfig(cfg)
	}

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func (cl *MqttClientSim) runSimulation(wg *sync.WaitGroup, clientNum int) {
	defer wg.Done()

	device := NewQuadsenseDevice(int32(clientNum), cl.UseSha)
	client, err := cl.newMqttClient(device.DeviceId)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(500)

	for i := 0; i != cl.NumRequests; i++ {
		if i != 0 {
			sleepJitter(cl.Frequency, 0.01)
		}
		for _, sensor := range device.Sensors {
			value := sensor.Value()
			if cl.UseCounter {
				value = fmt.Sprintf("%d,%s", i, value)
			}
			if cl.LogRequests {
				fmt.Fprintf(os.Stderr, "%s t: %s msg: %s\n", time.Now().Format("2006-01-02T15:04:05.999"), device.Topic(sensor), value)
			}
			sleepJitter(10*time.Millisecond, 0.1)
			client.Publish(device.Topic(sensor), cl.Qos, false, value)

		}
	}
}
