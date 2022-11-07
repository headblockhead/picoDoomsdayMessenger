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
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	// Setup an LED so that if there is an error, we know about it.
	led := machine.GPIO24
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	// Setup the RGB LED array.
	neopixelpin := machine.GP3
	neopixelpin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.New(neopixelpin)

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
	oldDeviceHighlightedItem := &picodoomsdaymessenger.MenuItemDefault

	// Store the last time that an LED animation frame was displayed.
	lastAnimationFrame := time.Now()

	// Store the last time that any of the buttons were pressed.
	lastButtonPress := time.Now()

	// Record the display size
	displayx, displayy := display.Size()

	// Set all the LEDs to black, do this multiple times to make sure they are all off.
	for i := 0; i < 8; i++ {
		err := displayLEDArray(&leds, [6]color.RGBA{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}})
		if err != nil {
			handleError(&display, &led, device, err)
		}
	}

	// Setup input reading. The columns are read and the rows are pulsed.
	buttonsCol1 := machine.GP4
	buttonsCol2 := machine.GP5
	buttonsCol3 := machine.GP6
	buttonsCol4 := machine.GP7
	buttonsCol5 := machine.GP8
	buttonsRow1 := machine.GP16
	buttonsRow2 := machine.GP17
	buttonsRow3 := machine.GP20
	buttonsRow4 := machine.SPI0_SDO_PIN
	buttonsRow5 := machine.GP22

	buttonsCol1.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol2.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol3.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol4.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol5.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsRow1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow4.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow5.Configure(machine.PinConfig{Mode: machine.PinOutput})

	buttonsRow1.Low()
	buttonsRow2.Low()
	buttonsRow3.Low()
	buttonsRow4.Low()
	buttonsRow5.Low()

	// Define the input maps
	buttons := [5][5]picodoomsdaymessenger.InputEvent{
		{picodoomsdaymessenger.InputEventNumber1, picodoomsdaymessenger.InputEventNumber2, picodoomsdaymessenger.InputEventNumber3, picodoomsdaymessenger.InputEventFunction1, picodoomsdaymessenger.InputEventUp},
		{picodoomsdaymessenger.InputEventNumber4, picodoomsdaymessenger.InputEventNumber5, picodoomsdaymessenger.InputEventNumber6, picodoomsdaymessenger.InputEventFunction2, picodoomsdaymessenger.InputEventDown},
		{picodoomsdaymessenger.InputEventNumber7, picodoomsdaymessenger.InputEventNumber8, picodoomsdaymessenger.InputEventNumber9, picodoomsdaymessenger.InputEventFunction3, picodoomsdaymessenger.InputEventLeft},
		{picodoomsdaymessenger.InputEventStar, picodoomsdaymessenger.InputEventNumber0, picodoomsdaymessenger.InputEventPound, picodoomsdaymessenger.InputEventFunction4, picodoomsdaymessenger.InputEventRight},
		{picodoomsdaymessenger.InputEventOpenMainMenu, picodoomsdaymessenger.InputEventOpenMessages, picodoomsdaymessenger.InputEventOpenPeople, picodoomsdaymessenger.InputEventOpenSettings, picodoomsdaymessenger.InputEventAccept},
	}

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

	// Main loop
	for {
		// Check the input if it has been long enough since the last button press.
		if lastButtonPress.Add(100 * time.Millisecond).Before(time.Now()) {
			buttonsRow1.High()
			col := checkInputCols(&buttonsCol1, &buttonsCol2, &buttonsCol3, &buttonsCol4, &buttonsCol5)
			if col != 0 {
				err := device.ProcessInputEvent(buttons[0][col-1])
				if err != nil {
					handleError(&display, &led, device, err)
					continue
				}
				lastButtonPress = time.Now()
			}
			buttonsRow1.Low()
			buttonsRow2.High()
			col = checkInputCols(&buttonsCol1, &buttonsCol2, &buttonsCol3, &buttonsCol4, &buttonsCol5)
			if col != 0 {
				err := device.ProcessInputEvent(buttons[1][col-1])
				if err != nil {
					handleError(&display, &led, device, err)
					continue
				}
				lastButtonPress = time.Now()
			}
			buttonsRow2.Low()
			buttonsRow3.High()
			col = checkInputCols(&buttonsCol1, &buttonsCol2, &buttonsCol3, &buttonsCol4, &buttonsCol5)
			if col != 0 {
				err := device.ProcessInputEvent(buttons[2][col-1])
				if err != nil {
					handleError(&display, &led, device, err)
					continue
				}
				lastButtonPress = time.Now()
			}
			buttonsRow3.Low()
			buttonsRow4.High()
			col = checkInputCols(&buttonsCol1, &buttonsCol2, &buttonsCol3, &buttonsCol4, &buttonsCol5)
			if col != 0 {
				err := device.ProcessInputEvent(buttons[3][col-1])
				if err != nil {
					handleError(&display, &led, device, err)
					continue
				}
				lastButtonPress = time.Now()
			}
			buttonsRow4.Low()
			buttonsRow5.High()
			col = checkInputCols(&buttonsCol1, &buttonsCol2, &buttonsCol3, &buttonsCol4, &buttonsCol5)
			if col != 0 {
				err := device.ProcessInputEvent(buttons[4][col-1])
				if err != nil {
					handleError(&display, &led, device, err)
					continue
				}
				lastButtonPress = time.Now()
			}
			buttonsRow5.Low()
		}
		// Update the display if the state changes
		if !reflect.DeepEqual(oldDeviceState, device.State) || !reflect.DeepEqual(oldDeviceHighlightedItem, device.State.HighlightedItem) {
			oldDeviceState = *device.State
			oldDeviceHighlightedItem = device.State.HighlightedItem
			frame, err := picodoomsdaymessenger.GetFrame(image.Rect(0, 0, int(displayx), int(displayy)), device)
			if err != nil {
				handleError(&display, &led, device, err)
				continue
			}
			err = displayImage(&display, frame)
			if err != nil {
				handleError(&display, &led, device, err)
				continue
			}
		}

		// Display the next animation frame if it has been long enough since the last frame.
		if lastAnimationFrame.Add(device.LEDAnimation.FrameDuration).Before(time.Now()) {
			device.LEDAnimation.CurrentFrame++
			if device.LEDAnimation.CurrentFrame >= len(device.LEDAnimation.Frames) {
				device.LEDAnimation.CurrentFrame = 0
			}
			displayLEDArray(&leds, device.LEDAnimation.Frames[device.LEDAnimation.CurrentFrame])
			lastAnimationFrame = time.Now()
		}

		time.Sleep(time.Millisecond * 1)
	}
}

func displayLEDArray(leds *ws2812.Device, ledlist [6]color.RGBA) error {
	for i := 0; i < len(ledlist); i++ {
		for j := 0; j < 8; j++ {
			err := leds.WriteByte(ledlist[i].G)
			if err != nil {
				return err
			}
			err = leds.WriteByte(ledlist[i].B)
			if err != nil {
				return err
			}
			err = leds.WriteByte(ledlist[i].R)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func checkInputCols(buttonsCol1, buttonsCol2, buttonsCol3, buttonsCol4, buttonsCol5 *machine.Pin) int {
	if buttonsCol1.Get() {
		return 1
	}
	if buttonsCol2.Get() {
		return 2
	}
	if buttonsCol3.Get() {
		return 3
	}
	if buttonsCol4.Get() {
		return 4
	}
	if buttonsCol5.Get() {
		return 5
	}
	return 0
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
