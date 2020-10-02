## run it

tinygo flash -target=microbit ./main.go

### GPIO breadboard

```
	SPI0_SCK_PIN  Pin = 23 // P13 on the board
	SPI0_MOSI_PIN Pin = 21 // P15 on the board
	SPI0_MISO_PIN Pin = 22 // P14 on the board

  const clock = machine.SPI0_SCK_PIN
	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       clock,
		MISO:      machine.SPI0_MISO_PIN,
		MOSI:      machine.SPI0_MOSI_PIN,
		Mode:      0})
```

### Gator Clips

```
	P0  Pin = 3
	P1  Pin = 2
	P2  Pin = 1

  const clock = machine.P1
	clock.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 4000000,
		SCK:       clock,
		MISO:      machine.P2,
		MOSI:      machine.P0,
		Mode:      0})
```

### colors

```
	bufOrange      = []byte{x, 0x00, 0x5a, 0xa8} // 0xd4, 0x51, 0x0b}
	bufYellow      = []byte{x, 0x08, 0x8c, 0xff} // 0xd4, 0x51, 0x0b}
	bufBlue        = []byte{x, 0xf5, 0x87, 0x42}
	bufGreen       = []byte{x, 0x0b, 0xd4, 0x4e}
	bufMag         = []byte{x, 0xd4, 0x0b, 0xa2}
```
