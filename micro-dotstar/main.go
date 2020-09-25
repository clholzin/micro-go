// Connects dotStar strip through SPI.
package main

import (
	"fmt"
	"machine"
	"time"
)

// cs is the pin used for Chip Select (CS). Change to whatever is in use on your board.
const clock = machine.SPI0_SCK_PIN

var (
	tx        []byte
	rx        []byte
	numpix    int  = 144
	x         byte = 0xef // brightness control E0 to FF
	r         byte = 0xff // red value
	g         byte = 0x0f // green value
	b         byte = 0xab // blue value
	buf            = []byte{x, b, g, r}
	bufOrange      = []byte{x, 0xd4, 0x51, 0x0b}
	bufBlue        = []byte{x, 0x42, 0x87, 0xf5}
	bufGreen       = []byte{x, 0x4e, 0xd4, 0x0b}
)

func main() {
	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMagRow := machine.LED_ROW_1
	ledMagRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMag := machine.LED_COL_1
	ledMag.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       clock,
		MISO:      machine.SPI0_MISO_PIN,
		MOSI:      machine.SPI0_MOSI_PIN,
		Mode:      0})

	ledMagRow.High()
	clock.High()
	var toggle int
	for {
		ledMag.Low()
		toggle++
		wleds := Leds(toggle)
		if toggle == 3 {
			toggle = 0
		}
		//rleds := make([]byte, len(wleds))
		clock.Low()
		err := machine.SPI0.Tx(wleds, nil)
		if err != nil {
			fmt.Printf("%s\n", err)
			clock.High()
			break
		}
		clock.High()
		//fmt.Printf("%v\n", rleds)
		ledMag.High()
		time.Sleep(300 * time.Millisecond)

	}
}

func Leds(toggle int) (tx []byte) {
	tx = append([]byte{}, []byte{0x00, 0x00, 0x00, 0x00}...)
	for i := 0; i < numpix; i++ {
		switch toggle {
		case 1:
			tx = append(tx, bufOrange...)
			break
		case 2:
			tx = append(tx, bufBlue...)
			break
		case 3:
			tx = append(tx, bufGreen...)
			break
		default:
			tx = append(tx, buf...)
		}
	}
	tx = append(tx, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
	return
}
