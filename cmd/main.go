package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	oaSim "github.com/openaircgn/openairMqttSimulator"
)

var (
	host            = flag.String("h", "tcp://localhost:1883", "broker url to connect to")
	numClients      = flag.Int("c", 1, "number of clients to simulate")
	frequency       = flag.Int("f", 10, "how many seconds to wait between sending measurements")
	numRequests     = flag.Int("n", 10, "number of total request, -1 for continuous")
	silent          = flag.Bool("s", false, "suppress detailed output about sent messages")
	useSha          = flag.Bool("sha", false, "use sha-1 hashed mac instead of raw mac")
	useCounter      = flag.Bool("counter", false, "add a counter to disambiguate messsages")
	qos             = flag.Int("qos", 0, "specify qos value (0,1,2 at most, at least, exactly once)")
	_version        = flag.Bool("version", false, "print version & exit")
	ca_pem          = flag.String("ca-pem", "", "path to a ca crt file (pem) to verify the server")
	client_cert_dir = flag.String("client-certs", ".", "path to search for client certs and keys, named client_id.pem and client_id.crt respectively")
	client_cfg      = flag.String("client-config", "", "path to client configuration, see test/client_list.json")

	version string
)

func banner() {
	fmt.Fprintf(os.Stderr, "version: %s\n", version)
}

func summary(list []oaSim.ClientCertConfig) {
	banner()
	fmt.Fprintf(os.Stderr, "simulating %d clients sending messages every %d seconds\n",
		*numClients, *frequency)

	if *numRequests > 0 {
		fmt.Fprintf(os.Stderr, "until each client has sent %d messages.\n", *numRequests)
	} else {
		fmt.Fprintf(os.Stderr, "continuously. Press Ctl-C to interrupt.\n")
	}

	fmt.Fprintf(os.Stderr, "target host: %s qos: %d\n\n", *host, *qos)

	fmt.Fprintf(os.Stderr, "client ids:\n")
	for _, cfg := range list {
		fmt.Fprintf(os.Stderr, "\t%s", cfg.ClientId)
		if cfg.MAC != "" {
			fmt.Fprintf(os.Stderr, "(mac:%s)", cfg.MAC)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}

	fmt.Fprintf(os.Stderr, "\n\n")

}
func createTLSConfig() *tls.Config {
	cfg := &tls.Config{}
	certPool := x509.NewCertPool()
	if *ca_pem != "" {
		if ca, err := ioutil.ReadFile(*ca_pem); err != nil {
			fmt.Fprintf(os.Stderr, "Could not read: %s (%v)\n", *ca_pem, err)
			flag.Usage()
			os.Exit(1)
		} else {
			certPool.AppendCertsFromPEM(ca)
		}

	}
	cfg.RootCAs = certPool
	return cfg
}

func populateClientCertConfigList() []oaSim.ClientCertConfig {
	var list []oaSim.ClientCertConfig
	for i := 0; i != *numClients; i++ {
		var clientId string
		if *useSha {
			clientId = oaSim.MakeMacSha(int32(i))
		} else {
			clientId = oaSim.MakeMAC(int32(i))
		}
		cfg := oaSim.ClientCertConfig{
			clientId,
			fmt.Sprintf("%s/%s.crt", *client_cert_dir, clientId),
			fmt.Sprintf("%s/%s.pem", *client_cert_dir, clientId),
			oaSim.MakeMAC(int32(i)),
		}

		list = append(list, cfg)
	}
	return list
}

func main() {
	flag.Parse()

	if *_version {
		banner()
		os.Exit(0)
	}

	var list []oaSim.ClientCertConfig
	var err error
	if *client_cfg != "" {
		list, err = oaSim.LoadClientCertConfigList(*client_cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read: %s (%v)\n", *client_cfg, err)
			flag.Usage()
			os.Exit(1)
		}
		*numClients = len(list)
	} else {
		list = populateClientCertConfigList()
	}

	if !*silent {
		summary(list)
	}

	tlsConfig := createTLSConfig()

	sim := oaSim.MqttClientSim{
		*host,
		tlsConfig,
		list,
		time.Duration(*frequency) * time.Second,
		*numRequests,
		!*silent,
		*useCounter,
		byte(*qos),
	}

	sim.RunSimulation()

}
