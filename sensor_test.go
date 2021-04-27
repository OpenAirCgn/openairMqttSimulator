package openairMqttSimulator

import (
	"fmt"
	"testing"
)

func TestIntGenerator(t *testing.T) {
	gen := IntGenerator{-100, 100, 3}
	for i := 0; i != 100; i++ {
		fmt.Println(gen.GenerateValue())
	}
}
func TestFloatGenerator(t *testing.T) {
	gen := FloatGenerator{-1.5, 50.5, 3}
	for i := 0; i != 100; i++ {
		fmt.Println(gen.GenerateValue())
	}
}
func TestBmeGenerator(t *testing.T) {
	for i := 0; i != 100; i++ {
		fmt.Println(bmeGenerator.GenerateValue())
	}
}
