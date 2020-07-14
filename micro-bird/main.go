package main

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/mag3110"
	"tinygo.org/x/drivers/mma8653"
)

var (
	i2c     = machine.I2C0
	buttonB = machine.Pin(machine.BUTTONB)
	pin0    = machine.Pin(machine.P0)
)

func main() {

	pin0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonB.Configure(machine.PinConfig{Mode: machine.PinInput})

	ledMagRow := machine.LED_ROW_1
	ledMagRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMag := machine.LED_COL_1
	ledMag.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledAccelRow := machine.LED_ROW_1
	ledAccelRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledAccel := machine.LED_COL_2
	ledAccel.Configure(machine.PinConfig{Mode: machine.PinOutput})

	i2c.Configure(machine.I2CConfig{SCL: machine.SCL_PIN, SDA: machine.SDA_PIN})

	var mag = mag3110.New(i2c)

	mag.Address = mag3110.Address
	mag.Configure()
	magConnected := mag.Connected()
	println("3-axis magnetometer Connected: ", magConnected)

	var accel = mma8653.New(i2c)

	accel.Address = mma8653.Address
	accel.Configure(mma8653.DataRate200Hz, mma8653.Sensitivity2G)
	accelConnected := accel.Connected()
	println("3-axis accelerometer Connected: ", accelConnected)
	acceler_mem := make([]uint, 3)
	mag_mem := make([]uint, 3)

	ledMagRow.High()
	ledAccelRow.High()
	pin0.Low()
	for {

		if !buttonB.Get() {
			println("buttonB")
			pin0.High()
			time.Sleep(250 * time.Millisecond)
			pin0.Low()
			continue
		}

		if magConnected {
			//func (d Device) ReadMagnetic() (x int16, y int16, z int16)
			x, y, z := mag.ReadMagnetic()
			println("mag: ", "X:", x, "Y:", y, "Z:", z)

			var diff = (mag_mem[0] / uint(x)) +
				(mag_mem[1] / uint(y)) +
				(mag_mem[2] / uint(z))
			println("magnetic diff: ", diff)
			if diff > uint(3) {
				ledMag.High()
				pin0.High()
			}

			mag_mem[0] = uint(x)
			mag_mem[1] = uint(y)
			mag_mem[2] = uint(z)
		}

		if accelConnected {
			//func (d Device) ReadAcceleration() (x int32, y int32, z int32, err error)
			x, y, z, err := accel.ReadAcceleration()
			if err == nil {
				println("accel: ", "X:", uint(x), "Y:", uint(y), "Z:", uint(z))

				var diff = (acceler_mem[0] / uint(x)) +
					(acceler_mem[1] / uint(y)) +
					(acceler_mem[2] / uint(z))

				if diff > uint(4) {
					ledAccel.High()
					pin0.High()
				}

				acceler_mem[0] = uint(x)
				acceler_mem[1] = uint(y)
				acceler_mem[2] = uint(z)
				println("acceleration diff: ", diff)
			}
		}

		time.Sleep(time.Millisecond * 250)
		ledMag.Low()
		ledAccel.Low()
		pin0.Low()

		//ledMagRow.Low()
		//ledAccelRow.Low()
		time.Sleep(time.Millisecond * 250)

	}

}
