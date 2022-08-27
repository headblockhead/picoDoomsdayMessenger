package main

import (
	"image"
	"image/color"
	"machine"
	"time"

	picodoomsdaymessenger "github.com/headblockhead/picoDoomsdayMessenger"
	"tinygo.org/x/drivers/ssd1306"
)

func main() {
	// Setup an LED so that if there is an error, we know about it.
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SDA:       machine.GP0,
		SCL:       machine.GP1,
	})
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.ClearDisplay()

	machine := picodoomsdaymessenger.NewMachine()

	displayx, displayy := display.Size()
	for {
		frame, err := picodoomsdaymessenger.GetFrame(image.Rect(0, 0, int(displayx), int(displayy)), machine)
		if err != nil {
			flashLED(&led, 1)
			return
		}
		err = displayImage(&display, frame)
		if err != nil {
			flashLED(&led, 2)
			return
		}
	}
}

func displayImage(display *ssd1306.Device, img image.Image) (err error) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			display.SetPixel(int16(x), int16(y), color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	err = display.Display()
	if err != nil {
		return err
	}
	return nil
}

func flashLED(led *machine.Pin, count int) {
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(300 * time.Millisecond)
		led.Low()
		time.Sleep(300 * time.Millisecond)
	}
}
