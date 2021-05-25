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
	BrokerUrl    string             // tcp://bla:1883 tls://... ws:// wss://
	TLSConfig    *tls.Config        // optional TLS configuration
	ClientConfig []ClientCertConfig // list of client ids (and optional tls certs) to use in simulation
	Frequency    time.Duration      // how long to sleep between "measurements"
	NumRequests  int                // total number of requests to send
	LogRequests  bool               // whether or not to be verbose
	UseCounter   bool               // use a counter as the first value transmitted
	Qos          byte
}

func (cl *MqttClientSim) RunSimulation() {
	var wg sync.WaitGroup
	for i := 0; i != len(cl.ClientConfig); i++ {
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

func setClientCertIfAvailable(cfg ClientCertConfig, tlsCfg *tls.Config) {

	cert, err := tls.LoadX509KeyPair(cfg.CRTFilename, cfg.PEMFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load client cert for %s (%v)\n", cfg.ClientId, err)
		return
	}
	tlsCfg.Certificates = []tls.Certificate{cert}
}

// creates a new (acutal, not simulated) mqtt client, connecting to the
// indicated URL using the provided client id.
func (cl *MqttClientSim) newMqttClient(clientNum int) (mqtt.Client, error) {
	options := mqtt.NewClientOptions()
	options.AddBroker(cl.BrokerUrl)
	options.SetClientID(cl.ClientConfig[clientNum].ClientId)

	if strings.HasPrefix(cl.BrokerUrl, "tls://") || strings.HasPrefix(cl.BrokerUrl, "wss://") {
		cfg := cl.TLSConfig
		if cfg == nil {
			cfg = &tls.Config{}
		}

		// attempt to set a client certificate if available
		setClientCertIfAvailable(cl.ClientConfig[clientNum], cfg)

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

	device := NewQuadsenseDevice(cl.ClientConfig[clientNum].ClientId)
	client, err := cl.newMqttClient(clientNum)
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
