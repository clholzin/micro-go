// Connects to a LIS3DH I2C accelerometer on the Adafruit Circuit Playground Express.
package main

import (
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/ws2812"
)

var i2c = machine.I2C1
var setTime = false

const (
	buttonA = machine.BUTTONA
	buttonB = machine.BUTTONB
)

func main() {

	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	buttonA.Configure(machine.PinConfig{Mode: machine.PinInput})
	buttonB.Configure(machine.PinConfig{Mode: machine.PinInput})

	ws := ws2812.New(neo)
	whiteleds := make([]color.RGBA, 10)
	blackleds := make([]color.RGBA, 10)

	for i := range whiteleds {
		whiteleds[i] = color.RGBA{R: 0xff, G: 0xff, B: 0xff}
	}

	for i := range blackleds {
		blackleds[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	}

	for {
		if buttonA.Get() {
			ws.WriteColors(whiteleds)
		}
		if buttonB.Get() {
			ws.WriteColors(blackleds)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
