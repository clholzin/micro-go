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
	tx     []byte
	rx     []byte
	numpix int = 4
)

func main() {
	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMagRow := machine.LED_ROW_1
	ledMagRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMag := machine.LED_COL_1
	ledMag.Configure(machine.PinConfig{Mode: machine.PinOutput})

	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 1000000, //4000000,
		SCK:       clock,
		MISO:      machine.SPI0_MISO_PIN,
		MOSI:      machine.SPI0_MOSI_PIN,
		Mode:      0})

	ledMagRow.High()
	clock.High()
	for {
		ledMag.Low()
		wleds := Leds()
		rleds := make([]byte, len(wleds))
		machine.SPI0.Tx(wleds, rleds)
		fmt.Printf("%v\n", rleds)
		ledMag.High()
		time.Sleep(100 * time.Millisecond)

	}
}

// Read analog data from channel
func Leds() (tx []byte) {
	tx = append([]byte{}, []byte{0x00, 0x00, 0x00, 0x00}...)
	for i := 0; i < numpix; i++ {
		tx = append(tx, []byte{0xFF, 0xFF, 0x00, 0x00}...)
	}
	tx = append(tx, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
	return
	//result = uint16((rx[1]&0x3))<<8 + uint16(rx[2])
}
