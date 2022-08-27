package main

import (
	"fmt"
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
	// Set the old machine state and old menu item to something that is not the starting value.
	oldMachineState := picodoomsdaymessenger.StateSettings
	oldDeviceMenuItem := 1

	// Record the display size
	displayx, displayy := display.Size()

	// Setup input reading
	clk := machine.GP4
	dta := machine.GP3
	sw := machine.GP2
	clk.Configure(machine.PinConfig{Mode: machine.PinInput})
	dta.Configure(machine.PinConfig{Mode: machine.PinInput})
	sw.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	var clkNow, clkPrv bool

	// Setup panic recovery
	defer func() {
		if err := recover(); err != nil {
			// If there is a panic, try to print details to the screen.
			frame, newErr := picodoomsdaymessenger.GetErrorFrame(image.Rect(0, 0, int(displayx), int(displayy)), device, fmt.Sprintf("%v", err))
			if newErr != nil {
				flashLED(&led, 1)
				return
			}
			newErr = displayImage(&display, frame)
			if newErr != nil {
				flashLED(&led, 2)
				return
			}
		}
	}()

	for {
		// Update the display if the state changes
		if oldMachineState != device.State || oldDeviceMenuItem != device.CurrentMenuItem {
			oldMachineState = device.State
			oldDeviceMenuItem = device.CurrentMenuItem
			frame, err := picodoomsdaymessenger.GetFrame(image.Rect(0, 0, int(displayx), int(displayy)), device)
			if err != nil {
				// If there is a error, try to print details to the screen.
				displayErrorFrame(&display, &led, device, err)
				return
			}
			err = displayImage(&display, frame)
			if err != nil {
				// If there is a error, try to print details to the screen.
				displayErrorFrame(&display, &led, device, err)
				return
			}
		}
		// Send a message to the device if an input is pressed

		// If the rotary encoder is pressed
		if !sw.Get() {
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventFire)
			if err != nil {
				// If there is a error, try to print details to the screen.
				displayErrorFrame(&display, &led, device, err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
		// If the rotary encoder is turned
		clkNow = clk.Get()
		if (clkNow != clkPrv) && clkNow {
			if dta.Get() {
				// Anti-Clockwise
				err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventLeft)
				if err != nil {
					// If there is a error, try to print details to the screen.
					displayErrorFrame(&display, &led, device, err)
					return
				}
			} else {
				// Clockwise
				err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventRight)
				if err != nil {
					// If there is a error, try to print details to the screen.
					displayErrorFrame(&display, &led, device, err)
					return
				}
			}
		}
		clkPrv = clkNow
		time.Sleep(time.Millisecond * 1)
	}
}

func displayErrorFrame(display *ssd1306.Device, led *machine.Pin, device *picodoomsdaymessenger.Device, inputerr error) (err error) {
	displayx, displayy := display.Size()
	frame, newErr := picodoomsdaymessenger.GetErrorFrame(image.Rect(0, 0, int(displayx), int(displayy)), device, inputerr.Error())
	if newErr != nil {
		flashLED(led, 1)
		return
	}
	newErr = displayImage(display, frame)
	if newErr != nil {
		flashLED(led, 2)
		return
	}
	return nil
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
