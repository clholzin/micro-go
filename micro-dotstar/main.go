// Connects dotStar strip through SPI.
package main

import (
	"encoding/hex"
	"fmt"
	"machine"
	"strconv"
	"time"
)

var (
	tx     []byte
	rx     []byte
	numpix int = 144

	x   byte = 0x46     // brightness control E0 to FF
	r   byte = 0xff     // red value
	g   byte = 0x0f     // green value
	b   byte = 0xab     // blue value
	buf      = bufWhite //[]byte{x, b, g, r}

	bufWhite   = []byte{x, 0xff, 0xff, 0xff}
	bufOrange  = []byte{x, 0x05, 0xa8, 0xff} //0x00, 0x5a, 0xa8} // 0xd4, 0x51, 0x0b}
	bufYellow  = []byte{x, 0x05, 0xea, 0xff} // 0xd4, 0x51, 0x0b}
	bufOrange2 = []byte{x, 0x00, 0xc4, 0xff} //[]byte{x, 0xf5, 0x87, 0x42}
	bufGreen   = []byte{x, 0x1e, 0xdb, 0x14}
	bufGreen2  = []byte{x, 0x1b, 0xde, 0x5f}
	bufMag     = []byte{x, 0xad, 0x8a, 0xff} // []byte{x, 0xd4, 0x0b, 0xa2}
	//0xff,0xce,0x0a
	colorStep  = 1
	bright     = 245
	greenStart = 113
	greenEnd   = 168
	red        = 255
	blue       = 5

	start, end int = 0, 10
)

const clock = machine.P1

func main() {

	ledMagRow := machine.LED_ROW_1
	ledMagRow.Configure(machine.PinConfig{Mode: machine.PinOutput})

	ledMag := machine.LED_COL_1
	ledMag.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// cs is the pin used for Chip Select (CS). Change to whatever is in use on your board.
	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       clock,
		MISO:      machine.P2,
		MOSI:      machine.P0,
		Mode:      0})

	ledMagRow.High()

	clock.High()
	//var toggle int
	for {
		ledMag.Low()

		/*wleds := append([]byte{}, []byte{0x00, 0x00, 0x00, 0x00}...)
		for i := 0; i < numpix; i++ {
			wleds = append(wleds, bufOrange...)
		}
		wleds = append(wleds, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
		clock.Low()
		err := machine.SPI0.Tx(wleds, nil)
		if err != nil {
			fmt.Printf("%s\n", err)
			clock.High()
			break
		}
		clock.High()
		time.Sleep(500 * time.Millisecond)*/
		down(clock)
		up(clock)

		ledMag.High()
		time.Sleep(100 * time.Millisecond)

	}
}

func up(clock machine.Pin) {
	for i := 0; i <= greenEnd-greenStart; i++ {
		wleds := pulse(greenStart + i)
		clock.Low()
		err := machine.SPI0.Tx(wleds, nil)
		if err != nil {
			fmt.Printf("%s\n", err)
			clock.High()
			break
		}
		clock.High()
		time.Sleep(50 * time.Millisecond)
	}
}

func down(clock machine.Pin) {
	for i := 0; i <= greenEnd-greenStart; i++ {
		wleds := pulse(greenEnd - i)
		clock.Low()
		err := machine.SPI0.Tx(wleds, nil)
		if err != nil {
			fmt.Printf("%s\n", err)
			clock.High()
			break
		}
		clock.High()
		time.Sleep(50 * time.Millisecond)
	}
}

func pulse(green int) (tx []byte) {
	tx = append(tx, []byte{0x00, 0x00, 0x00, 0x00}...)
	for i := 0; i < numpix; i++ {
		tx = append(tx, color(bright, red, green, blue)...)
	}
	tx = append(tx, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
	return
}

func color(bright, red, green, blue int) (c []byte) {
	xx, _ := hex.DecodeString(strconv.FormatUint(uint64(bright), 16))
	//rr, _ := hex.DecodeString(strconv.FormatUint(uint64(red), 16))
	gg, _ := hex.DecodeString(strconv.FormatUint(uint64(green), 16))
	//bb, _ := hex.DecodeString(strconv.FormatUint(uint64(blue), 16))
	c = append(c, xx...)
	c = append(c, []byte{0x05}...)
	c = append(c, gg...)
	c = append(c, []byte{0xff}...)
	return
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
			tx = append(tx, bufWhite...)
		}
	}
	tx = append(tx, []byte{0xFF, 0xFF, 0xFF, 0xFF}...)
	return

	/*wleds := All(toggle)
	toggle++
	if toggle > 6 {
		toggle = 0
	}*/
	//rleds := make([]byte, len(wleds))
}
