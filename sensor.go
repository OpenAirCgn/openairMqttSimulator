package openairMqttSimulator

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

func MakeMAC(id int32) string {
	bs := make([]byte, 6)
	binary.PutVarint(bs, int64(id))
	for i, _ := range bs {
		bs[i] = byte((id >> (i * 8)) & 0xff)
	}
	return net.HardwareAddr(bs).String()
}

func MakeMacSha(id int32) string {

	bs := make([]byte, 6)
	binary.PutVarint(bs, int64(id))
	for i, _ := range bs {
		bs[i] = byte((id >> (i * 8)) & 0xff)
	}
	return fmt.Sprintf("%x", sha1.Sum(bs))
}

type Device struct {
	DeviceId string
	Sensors  []Sensor
}

func NewQuadsenseDevice(id int32, useSha bool) Device {
	var deviceId string
	if useSha {
		deviceId = MakeMacSha(id)
	} else {
		deviceId = MakeMAC(id)
	}
	device := Device{
		deviceId,
		[]Sensor{NO2, NO, OX, CO, BME0, BME1, SDS011},
	}
	return device
}

func (dev *Device) Topic(s Sensor) string {
	return fmt.Sprintf("/%s/%s", dev.DeviceId, s.Name)
}

// Sensor contains the name, included in the topic as well as the
// Generator to create random values.
type Sensor struct {
	Name      string
	Generator Generator
}

func (s *Sensor) Value() string {
	return s.Generator.GenerateValue()
}

// Interface to produce random values Different implemenations to
// produce values appropriate for different sensors being simulated.
// Current impls are stateless, but one can imagine gradually changing
// values or replay for more realistic test scenarios.
type Generator interface {
	GenerateValue() string
}

var alphaGenerator = &IntGenerator{-300, 3000, 2}

var NO2 = Sensor{
	"NO2-B43F",
	alphaGenerator,
}
var NO = Sensor{
	"NO-B4",
	alphaGenerator,
}
var OX = Sensor{
	"OX-B431",
	alphaGenerator,
}
var CO = Sensor{
	"CO-B4",
	alphaGenerator,
}

var bmeGenerator = &MultiGen{
	[]Generator{
		&FloatGenerator{-20.0, 45.0, 1},
		&FloatGenerator{10.0, 90.0, 1},
		&FloatGenerator{950.0, 1050.0, 1},
	},
}
var BME0 = Sensor{
	"BME280-0",
	bmeGenerator,
}

var BME1 = Sensor{
	"BME280-1",
	bmeGenerator,
}

var SDS011 = Sensor{
	"SDS011",
	&FloatGenerator{0.0, 100.0, 2},
}

type IntGenerator struct {
	Min   int
	Max   int
	Count int // how many csv values in the range
}

func (g *IntGenerator) GenerateValue() string {
	str := ""
	for i := 0; i != g.Count; i++ {
		if i != 0 {
			str += ","
		}
		value := rand.Intn(g.Max - g.Min)
		str += fmt.Sprintf("%d", value+g.Min)
	}
	return str

}

type FloatGenerator struct {
	Min   float32
	Max   float32
	Count int
}

func (g *FloatGenerator) GenerateValue() string {
	str := ""
	for i := 0; i != g.Count; i++ {
		if i != 0 {
			str += ","
		}
		value := rand.Float32()
		str += fmt.Sprintf("%0.2f", g.Min+(value*(g.Max-g.Min)))
	}
	return str

}

type MultiGen struct {
	Generators []Generator
}

func (g *MultiGen) GenerateValue() string {
	str := ""
	for i, g := range g.Generators {
		if i != 0 {
			str += ","
		}
		str += fmt.Sprintf("%s", g.GenerateValue())
	}
	return str
}
