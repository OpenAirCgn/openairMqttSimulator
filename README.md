# Introduction

This repository contains a commandline client meant to simulate the MQTT
behavior of the Open Air device.

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

