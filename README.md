# Introduction

This repository contains a commandline client meant to simulate the MQTT
behavior of the Open Air device.

# Usage

CLI binaries for Windows, Mac and Linux are available in from the Github
releases tab.

	Usage of openair_mqtt_sim:
	  -c int
		number of clients to simulate (default 1)
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

# Building

Source the `xcompile.sh` script to build binaries for Windows, Mac and Linux.

# Future Directions

Code will be expanded to handle TLS connections as well as TLS client
authentication.

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


