// Connects to a LIS3DH I2C accelerometer on the Adafruit Circuit Playground Express.
package main

import (
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/lis3dh"
	"tinygo.org/x/drivers/ws2812"
)

var i2c = machine.I2C1
var setTime = false

func main() {
	i2c.Configure(machine.I2CConfig{SCL: machine.SCL1_PIN, SDA: machine.SDA1_PIN})

	accel := lis3dh.New(i2c)
	accel.Address = lis3dh.Address1 // address on the Circuit Playground Express
	accel.Configure()
	accel.SetRange(lis3dh.RANGE_2_G)

	println(accel.Connected())

	neo := machine.NEOPIXELS
	neo.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ws := ws2812.New(neo)
	whiteleds := make([]color.RGBA, 10)
	blackleds := make([]color.RGBA, 10)

	for i := range whiteleds {
		whiteleds[i] = color.RGBA{R: 0xff, G: 0xff, B: 0xff}
	}

	for i := range blackleds {
		blackleds[i] = color.RGBA{R: 0x00, G: 0x00, B: 0x00}
	}

	acceler_mem := make([]uint, 3)
	for {
		x, y, z, _ := accel.ReadAcceleration()
		println("X:", uint(x), "Y:", uint(y), "Z:", uint(z))

		var diff = (acceler_mem[0] / uint(x)) +
			(acceler_mem[1] / uint(y)) +
			(acceler_mem[2] / uint(z))

		acceler_mem[0] = uint(x)
		acceler_mem[1] = uint(y)
		acceler_mem[2] = uint(z)

		println("dff ", uint(diff))
		if uint(diff) > 5 {
			ws.WriteColors(whiteleds)
			time.Sleep(time.Minute * 10) //time.Sleep(time.Minute)
		} else {
			ws.WriteColors(blackleds)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
