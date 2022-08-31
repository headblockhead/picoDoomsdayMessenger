package main

import (
	"fmt"
	"image"
	"image/color"
	"machine"
	"reflect"
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
	oldDeviceState := picodoomsdaymessenger.StateDefault
	oldDeviceHighlightedItem := picodoomsdaymessenger.GlobalMenuItemDefault

	// Record the display size
	displayx, displayy := display.Size()

	// Setup input reading
	controlUpSwitch := machine.GP4
	controlDownSwitch := machine.GP3
	controlConfirmSwitch := machine.GP2
	controlUpSwitch.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	controlDownSwitch.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	controlConfirmSwitch.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// Panic recovery
	defer func() {
		if err := recover(); err != nil {
			// Communicate that an error happened.
			flashLED(&led, 1, 300)
			// If there is a panic, try to print details to the screen.
			// The handleError() function cannot be used here as it requires an error.
			// err in this case is not an error but an interface.
			// So we use fmt.Sprintf("%v", err) to write details to the screen.
			frame, newErr := picodoomsdaymessenger.GetErrorFrame(image.Rect(0, 0, int(displayx), int(displayy)), device, fmt.Sprintf("%v", err))
			if newErr != nil {
				flashLED(&led, 2, 300)
				return
			}
			newErr = displayImage(&display, frame)
			if newErr != nil {
				flashLED(&led, 2, 300)
				return
			}
		}
	}()

	for {
		if !controlConfirmSwitch.Get() {
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventFire)
			if err != nil {
				handleError(&display, &led, device, err)
				return
			}
			time.Sleep(time.Millisecond * 250)
		}
		if !controlUpSwitch.Get() {
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventUp)
			if err != nil {
				handleError(&display, &led, device, err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
		if !controlDownSwitch.Get() {
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventDown)
			if err != nil {
				handleError(&display, &led, device, err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
		// Update the display if the state changes
		if !reflect.DeepEqual(oldDeviceState, device.State) || !reflect.DeepEqual(oldDeviceHighlightedItem, device.State.HighlightedItem) {
			oldDeviceState = *device.State
			oldDeviceHighlightedItem = device.State.HighlightedItem
			frame, err := picodoomsdaymessenger.GetFrame(image.Rect(0, 0, int(displayx), int(displayy)), device)
			if err != nil {
				handleError(&display, &led, device, err)
				return
			}
			err = displayImage(&display, frame)
			if err != nil {
				handleError(&display, &led, device, err)
				return
			}
		}
		time.Sleep(time.Millisecond * 1)
	}
}

// handleError takes in an error and communicates it to the user.
func handleError(display *ssd1306.Device, led *machine.Pin, device *picodoomsdaymessenger.Device, inputerr error) {
	// Communicate that an error happened.
	flashLED(led, 1, 300)
	// Try to get details to print to the screen
	displayx, displayy := display.Size()
	frame, newErr := picodoomsdaymessenger.GetErrorFrame(image.Rect(0, 0, int(displayx), int(displayy)), device, inputerr.Error())
	if newErr != nil {
		// If we can't do that, resort to signaling with the LED
		flashLED(led, 2, 300)
		return
	}
	// Try to print the details to the screen
	newErr = displayImage(display, frame)
	if newErr != nil {
		// If we can't do that either, resort to signaling with the LED
		flashLED(led, 2, 300)
		return
	}
}

// displayImage takes in an image and writes it to the screen.
func displayImage(display *ssd1306.Device, img image.Image) (err error) {
	// Put the image into the buffer.
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			display.SetPixel(int16(x), int16(y), color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	// Show the buffer.
	err = display.Display()
	if err != nil {
		return err
	}
	return nil
}

// flashLED will toggle an LED a certain amount of times and will wait a certain amount of time between toggles.
func flashLED(led *machine.Pin, count int, delay time.Duration) {
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(delay * time.Millisecond)
		led.Low()
		time.Sleep(delay * time.Millisecond)
	}
}
