// Connects dotStar strip through SPI.
package main

import (
	"fmt"
	"machine"
	"time"
)

var (
	tx     []byte
	rx     []byte
	numpix int = 144

	x   byte = 0xef     // brightness control E0 to FF
	r   byte = 0xff     // red value
	g   byte = 0x0f     // green value
	b   byte = 0xab     // blue value
	buf      = bufWhite //[]byte{x, b, g, r}

	bufWhite   = []byte{x, 0xff, 0xff, 0xff}
	bufOrange  = []byte{x, 0x00, 0x5a, 0xa8} // 0xd4, 0x51, 0x0b}
	bufYellow  = []byte{x, 0xff, 0xea, 0x05} // 0xd4, 0x51, 0x0b}
	bufOrange2 = []byte{x, 0x00, 0xc4, 0xff} //[]byte{x, 0xf5, 0x87, 0x42}
	bufGreen   = []byte{x, 0x1e, 0xdb, 0x14}
	bufGreen2  = []byte{x, 0x1b, 0xde, 0x5f}
	bufMag     = []byte{x, 0xad, 0x8a, 0xff} // []byte{x, 0xd4, 0x0b, 0xa2}

	start, end int = 0, 10
)

func main() {

	ledMagRow := machine.LED_ROW_1
	ledMagRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMag := machine.LED_COL_1
	ledMag.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// cs is the pin used for Chip Select (CS). Change to whatever is in use on your board.
	const clock = machine.P1
	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       clock,
		MISO:      machine.P2,
		MOSI:      machine.P0,
		Mode:      0})

	ledMagRow.High()

	clock.High()
	var toggle int
	for {
		ledMag.Low()
		wleds := All(toggle)
		toggle++
		if toggle > 2 {
			toggle = 2
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

		/*for i := 0; i < numpix; i++ {
			copy := wleds[:]
			en := i + 1
			if en > numpix {
				en = numpix
			}
			copy = Zip(i, en, copy)
			clock.Low()
			err := machine.SPI0.Tx(copy, nil)
			if err != nil {
				fmt.Printf("%s\n", err)
				clock.High()
				break
			}
			clock.High()
			time.Sleep(50 * time.Millisecond)

		}*/

		//fmt.Printf("%v\n", rleds)
		ledMag.High()
		time.Sleep(1000 * time.Millisecond)

	}
}

func Zip(start, end int, tx []byte) (tz []byte) {
	startbyte := tx[:start*4]
	endbyte := tx[end*4:]
	runbyte := []byte{}
	for i := 0; i < (start - end); i++ {
		runbyte = append(runbyte, buf...)
	}
	tz = append(startbyte, runbyte...)
	tz = append(tz, endbyte...)
	return tz
}

func All(toggle int) (tx []byte) {
	tx = append(tx, []byte{0x00, 0x00, 0x00, 0x00}...)
	for i := 0; i < numpix; i++ {
		switch toggle {
		case 1:
			tx = append(tx, bufOrange...)
			break
		case 2:
			tx = append(tx, bufYellow...)
			break
		case 3:
			tx = append(tx, bufGreen...)
			break
		case 4:
			tx = append(tx, bufGreen2...)
			break
		case 5:
			tx = append(tx, bufOrange2...)
		case 6:
			tx = append(tx, bufMag...)
		default:
			tx = append(tx, buf...)
		}
	}
	tx = append(tx, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
	return
}
