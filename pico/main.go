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

	// Setup the display
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

	// Create a new Machine
	device := picodoomsdaymessenger.NewDevice()
	oldMachineState := picodoomsdaymessenger.StateSettings

	// Record the display size
	displayx, displayy := display.Size()

	// Setup input reading
	clk := machine.GP4
	dta := machine.GP3
	sw := machine.GP2

	clk.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	dta.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	sw.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	var clkNow, clkPrv bool

	for {
		// Update the display if the state changes
		if oldMachineState != device.State {
			frame, err := picodoomsdaymessenger.GetFrame(image.Rect(0, 0, int(displayx), int(displayy)), device)
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
		// Send a message to the device if an input is pressed

		// If the rotary encoder is pressed
		if !sw.Get() {
			flashLED(&led, 1)
			defer func() {
				if err := recover(); err != nil {
					flashLED(&led, 10)
				}
			}()
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventFire)
			if err != nil {
				flashLED(&led, 3)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
		// If the rotary encoder is turned
		clkNow = clk.Get()
		if (clkNow != clkPrv) && clkNow {
			if dta.Get() {
				flashLED(&led, 1)
				// Anti-Clockwise
				err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventLeft)
				if err != nil {
					flashLED(&led, 5)
					return
				}
			} else {
				flashLED(&led, 2)
				// Clockwise
				err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventRight)
				if err != nil {
					flashLED(&led, 4)
					return
				}
			}
		}
		clkPrv = clkNow
		time.Sleep(time.Millisecond * 1)
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
