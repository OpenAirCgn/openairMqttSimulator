# Introduction

This repository contains a commandline client meant to simulate the MQTT
behavior of the Open Air device.

# Usage

CLI binaries for Windows, Mac and Linux are available in from the Github
releases tab.

	Usage of openair_mqtt_sim:
	  -c int
		number of clients to simulate (default 1)
	  -ca-pem string
		path to a ca crt file (pem) to verify the server cert
	  -client-certs string
		path to search for client certs and keys, named client_id.pem and client_id.crt respectively (default ".")
          -client-config string
    	        path to client configuration, see set "Client Configuration", below
	  -counter
		add a counter to disambiguate messsages
	  -f int
		how many seconds to wait between sending measurements (default 10)
	  -h string
		broker url to connect to (default "tcp://localhost:1883")
	  -n int
		number of total request, -1 for continuous (default 10)
	  -qos int
		specify qos value (0,1,2 at most, at least, exactly once)
	  -s	suppress detailed output about sent messages
	  -sha
		use sha-1 hashed mac instead of raw mac
	  -version
		print version & exit


E.g. in order to simulate 5 clients sending messages each every 15
minutes a total of 10 times each:

	$ openair_mqtt_sim -c 5 -n 10 -f 900 
					# 15min = 15 * 60 seconds

Use the `-sha` flag f a sha-1 hash of the mac should be used for both
the clientid and within the topic.

QOS (quality of service) is set to 0 (at most once) by default. To
experiment with other settings, use the `-qos` flag. Similarly use
`-counter` to prefix each set of values with a per client counter.

# Client Configuration

By default, you specify the number of clients to use and the simulator
"generates" them, i.e. a series of `ClientId`s corresponding to MAC
addresses `00:00:00:00:00:00`, `01:00:00:00:00:00`, etc. are generated.
If the `-sha` flag is present, this value is hashed.

If the simulation finds corresponding `crt` and `pem` files in the
`-client-certs-dir` directory, these certificates are used for TLS
client authentication.

An alternative approach is to provide a configuration file describing
the `ClientId`s and certificate files to be used. The location of this
file is specified using the `-client-config` flag, an example file is
located [here](test/client_list.json)


# Preliminary TLS support

The simulation will attempt to use TLS in case the broker url is prefixed with
`tls://` or `wss://` instead of `tcp://` and `ws://`.

In case a CA certificate needs to be provided for self-signed server
certificates, it may be passed with the `-ca-pem` flag. In case client TLS
authentication should be used, the path containing client certificates
and key files should be provided with the `-client-certs` flag.
Certificates and key files corresponding to a particular client ID
should be named `<clientID>.crt` and `<clientID>.pem` respectively.

Currently, pem files may NOT been protected by passphrase. See #Testing

# Building

Source the `xcompile.sh` script to build binaries for Windows, Mac and Linux.

# Testing

This section describes a simple test setup using a public MQTT server: test.mosquitto.org 1883

Subscribe using e.g: `mosquitto_sub -h test.mosquitto.org -p 1883 -t
'/7722745105e9e02e8f1aaf17f7b3aac5c56cd805/#'` to monitor transfer.

In case TLS is used with the test.mosquitto.org server, connect to port
8883 and use a `tls://` url, i.e. `tls://test.mosquitto.org:8883` instead of
`tcp://test.mosquitto.org:1883`

TLS *client authentication* my be tested against port 8884. Ten client
certificates have been prepared and signed by mosquitto. These are
available in the `test` directory. The client IDs correspond to the
hashed version of the simulated MACs:

    00:00:00:00:00:00
    ...
    09:00:00:00:00:00

A sample commandline to test client authentication would look like:

    $ openair_mqtt_sim \
          -h tls://test.mosquitto.org:8884 \
          -sha    # use hashed macs (sample crt's assume hashes are being used) \
          -c 10   # use all ten clients we have certs for \
          -ca-pem test/mosquitto.org.crt \
                  # mosquitto uses a self signed cert, need to provide CA cert to verify \
          -client_certs test # directory containing client certs 

The `test` directory also contains scripts detailing how the test cert
material was generated using the openssl cli. CSR's were signed using:
https://test.mosquitto.org/ssl/

# Future Directions

Auto-create self-signed client certificates for local testing.

# MQTT Behavior

## Authentication

Typically the MAC address of the ESP32 module serves as the client id of
the device. Optionally, the device may use TLS client certification to
verify its identity. In this case, the CN entry in the client
certificate contains the client id.

Optionally, a SHA hash of the MAC may be used in case of privacy
concerns. 

For the sake of this simulation, randomly generated mac addresses are
generated for each client.


## Topic structure

	topic      := “/” <device_id> ”/” <sensor>
	device_id  := optional_hash(device_MAC)
	sensor     := “NO2-B43F” | “NO-B4” | “XS-B431” | “CO-B4” | “BME280-0” | “BME280-1” | "SDS011" Feinstaub (PM25/PM10)
	counter    := <incremented once per send>
	values     := <value_alpha> | <value_bme> | <value_dust>
	value_alpha:= working_electrode_value “,” aux_electrode_value
	value_bme  := tmp “,” hum “,” pressure 
	value_dust := <pm25> "," <pm10>
	we_value   := value
	ae_value   := value
	temp       := value
	hum        := value
	pressure   := value
	pm_10      := value     
	pm_25      := value
	value := ASCII Digits, ".",  preceded with an optional Sign ("-")


- we_value ADC reading of the working electrode
- ae_value ADC reading of the auxiliary electrode

The "raw" form of sensor value is transmitted to allow for centralised
compensation of values. The concrete semantics of "raw" values will be
specified per sensor if necessary.

Groups of values originating from a single sensor will typically be
transmitted to a single topic, seperated by comma. E.g. The BME280 is
regarded as a single sensor which contributes three measurements (temp,
humidity and air pressure) and not as a thermometer, hygrometer and
barometer each with a single measurement. The is regarded as a rule of
thumb and is flexibel for change depending on sensor characteristics.


