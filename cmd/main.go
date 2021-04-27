package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	oaSim "github.com/openaircgn/openairMqttSimulator"
)

var (
	host        = flag.String("h", "tcp://localhost:1883", "broker url to connect to")
	numClients  = flag.Int("c", 1, "number of clients to simulate")
	frequency   = flag.Int("f", 10, "how many seconds to wait between sending measurements")
	numRequests = flag.Int("n", 10, "number of total request, -1 for continuous")
	silent      = flag.Bool("s", false, "suppress detailed output about sent messages")
	useSha      = flag.Bool("sha", false, "use sha-1 hashed mac instead of raw mac")
)

func summary() {
	fmt.Fprintf(os.Stderr, "simulating %d clients sending messages every %d seconds\n",
		*numClients, *frequency)
	if *numRequests > 0 {
		fmt.Fprintf(os.Stderr, "until each client has sent %d messages.\n", *numRequests)
	} else {
		fmt.Fprintf(os.Stderr, "continuously. Press Ctl-C to interrupt.\n")
	}
	fmt.Fprintf(os.Stderr, "target host: %s\n\n", *host)

	fmt.Fprintf(os.Stderr, "client ids:\n")
	for i := 0; i != *numClients; i++ {
		id := oaSim.MakeMAC(int32(i))
		if *useSha {
			shaId := oaSim.MakeMacSha(int32(i))
			fmt.Fprintf(os.Stderr, "\t%s (mac: %s)\n", shaId, id)
		} else {
			fmt.Fprintf(os.Stderr, "\t%s\n", id)
		}
	}

	fmt.Fprintf(os.Stderr, "\n\n")

}

func main() {
	flag.Parse()

	if !*silent {
		summary()
	}

	sim := oaSim.MqttClientSim{
		*host,
		*numClients,
		time.Duration(*frequency) * time.Second,
		*numRequests,
		!*silent,
		*useSha,
	}

	sim.RunSimulation()

}
