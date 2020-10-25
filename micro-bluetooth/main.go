package main

import (
	"machine"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

var (
	hasColorChange bool = false
	brightness     byte = 0x50 // brightness control E0 to FF
	ledColor            = []byte{brightness, 0xff, 0xff, 0xff}
	numpix         int  = 144
	serviceUUID         = [16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
	charUUID            = [16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
)

const clock = machine.P1

func main() {

	ledMagRow := machine.LED_ROW_1
	ledMagRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMag := machine.LED_COL_1
	ledMag.Configure(machine.PinConfig{Mode: machine.PinOutput})

	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       clock,
		SDO:       machine.P0,
		SDI:       machine.P2,
		Mode:      0})

	println("starting")
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "LED colors",
	}))
	must("start adv", adv.Start())

	var ledColorCharacteristic bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.NewUUID(serviceUUID),
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &ledColorCharacteristic,
				UUID:   bluetooth.NewUUID(charUUID),
				Value:  ledColor[:],
				Flags:  bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicWritePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					if offset != 0 || len(value) != 3 {
						return
					}
					println(value[0], value[1], value[2])
					ledColor[1] = value[0]
					ledColor[2] = value[1]
					ledColor[3] = value[2]
					hasColorChange = true
				},
			},
		},
	}))

	run := func() {
		ledMag.Low()
		wleds := append([]byte{}, []byte{0x00, 0x00, 0x00, 0x00}...)
		for i := 0; i < numpix; i++ {
			wleds = append(wleds, ledColor...)
		}
		wleds = append(wleds, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
		clock.Low()
		err := machine.SPI0.Tx(wleds, nil)
		if err != nil {
			println(err)
		}
		ledMag.High()
		clock.High()
	}

	ledMagRow.High()
	clock.High()

	run()

	for {

		if hasColorChange {
			println("starting led write")
			hasColorChange = false
			run()
			println("finishing led write")
		}

		time.Sleep(100 * time.Millisecond)

	}
}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
