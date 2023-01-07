package main

import (
	"fmt"
	"image"
	"image/color"
	"machine"
	"reflect"
	"time"

	picodoomsdaymessenger "github.com/headblockhead/picoDoomsdayMessenger"
	"github.com/headblockhead/tinygorfm9x"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/drivers/ws2812"
)

func main() {
	time.Sleep(2 * time.Second) // Wait for the USB serial to be ready.
	// Setup an LED so that if there is an error, we know about it.
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led.Low()

	// Setup the display
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.TWI_FREQ_400KHZ,
		SDA:       machine.GPIO0,
		SCL:       machine.GPIO1,
	})
	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})
	display.ClearDisplay()

	// Record the display size
	displayx, displayy := display.Size()

	// Create a new Machine
	device, err := picodoomsdaymessenger.NewDevice()
	if err != nil {
		handleError(&display, &led, device, err)
	}

	// Set the old machine state and old menu item to something that is not the starting value.
	oldDeviceState := picodoomsdaymessenger.StateDefault
	oldDeviceHighlightedItemIndex := 0

	// Set up panic recovery
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

	// Setup the RGB LED array.
	neopixelpin := machine.D6
	neopixelpin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.New(neopixelpin)

	// Clear the LED array.
	err = displayLEDArray(&leds, [6]color.RGBA{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}})
	if err != nil {
		handleError(&display, &led, device, err)
	}

	// Store the last time that an LED animation frame was displayed.
	lastAnimationFrame := time.Now()

	// Store the last time that any of the buttons were pressed.
	lastButtonPress := time.Now()

	// Setup the RFM9x radio.
	rfm := tinygorfm9x.RFM9x{
		SPIDevice: *machine.SPI1,
	}
	err = rfm.Init(tinygorfm9x.Options{
		FrequencyMHz:      868,
		ResetPin:          machine.LORA_RESET,
		CSPin:             machine.LORA_CS,
		DIO0Pin:           machine.LORA_DIO0,
		DIO1Pin:           machine.LORA_DIO1,
		DIO2Pin:           machine.LORA_DIO2,
		EnableCRCChecking: true,
	})
	if err != nil {
		handleError(&display, &led, device, err)
	}

	err = rfm.StartReceive()
	if err != nil {
		handleError(&display, &led, device, err)
	}

	rfm.OnReceivedPacket = func(packet tinygorfm9x.Packet) {
		err = device.ReceiveFromRadio(packet.Payload)
		if err != nil {
			handleError(&display, &led, device, err)
		}
	}

	device.SendUsingRadio = func(packet []byte) (err error) {
		println("Sending packet: " + string(packet))
		err = rfm.Send(packet)
		if err != nil {
			return err
		}
		println("Done packet: " + string(packet))
		return err
	}

	c := device.NewConversation(picodoomsdaymessenger.PersonYou)
	c.Messages = append(c.Messages, picodoomsdaymessenger.Message{
		Person: picodoomsdaymessenger.PersonYou,
		Text:   "Hello, world!",
	})
	c.Name = "New Message"
	c.HighlightedMessageIndex = 0
	device.UpdateConversationsMenu()

	// Setup input reading. The columns are read and the rows are pulsed.
	buttonsCol1 := machine.D9
	buttonsCol1.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol2 := machine.D10
	buttonsCol2.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol3 := machine.D11
	buttonsCol3.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol4 := machine.D12
	buttonsCol4.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	buttonsCol5 := machine.D13
	buttonsCol5.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	buttonsRow1 := machine.GPIO16
	buttonsRow1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow1.Low()
	buttonsRow2 := machine.GPIO17
	buttonsRow2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow2.Low()
	buttonsRow3 := machine.GPIO20
	buttonsRow3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow3.Low()
	buttonsRow4 := machine.GPIO23
	buttonsRow4.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow4.Low()
	buttonsRow5 := machine.GPIO22
	buttonsRow5.Configure(machine.PinConfig{Mode: machine.PinOutput})
	buttonsRow5.Low()

	// Define the locations of the buttons in the button array.
	buttons := [5][5]picodoomsdaymessenger.InputEvent{
		{picodoomsdaymessenger.InputEventNumber1, picodoomsdaymessenger.InputEventNumber2, picodoomsdaymessenger.InputEventNumber3, picodoomsdaymessenger.InputEventFunction1, picodoomsdaymessenger.InputEventUp},
		{picodoomsdaymessenger.InputEventNumber4, picodoomsdaymessenger.InputEventNumber5, picodoomsdaymessenger.InputEventNumber6, picodoomsdaymessenger.InputEventFunction2, picodoomsdaymessenger.InputEventDown},
		{picodoomsdaymessenger.InputEventNumber7, picodoomsdaymessenger.InputEventNumber8, picodoomsdaymessenger.InputEventNumber9, picodoomsdaymessenger.InputEventFunction3, picodoomsdaymessenger.InputEventLeft},
		{picodoomsdaymessenger.InputEventStar, picodoomsdaymessenger.InputEventNumber0, picodoomsdaymessenger.InputEventPound, picodoomsdaymessenger.InputEventFunction4, picodoomsdaymessenger.InputEventRight},
		{picodoomsdaymessenger.InputEventOpenMainMenu, picodoomsdaymessenger.InputEventOpenConversations, picodoomsdaymessenger.InputEventOpenPeople, picodoomsdaymessenger.InputEventOpenSettings, picodoomsdaymessenger.InputEventAccept},
	}

	// Main program loop.
	for {
		// Check the input if it has been long enough since the last button press.
		if lastButtonPress.Add(200 * time.Millisecond).Before(time.Now()) {
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

		// Update the display if the state has changed.
		if !reflect.DeepEqual(oldDeviceState, device.State) || !(oldDeviceHighlightedItemIndex == device.State.HighlightedItemIndex) {
			oldDeviceState = *device.State
			oldDeviceHighlightedItemIndex = device.State.HighlightedItemIndex
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
			if device.LEDAnimation.CurrentFrame >= len(device.LEDAnimation.Frames) {
				device.LEDAnimation.CurrentFrame = 0
			}
			displayLEDArray(&leds, device.LEDAnimation.Frames[device.LEDAnimation.CurrentFrame])
			device.LEDAnimation.CurrentFrame++
			lastAnimationFrame = time.Now()
		}
	}
}

// displayLEDArray displays the given RGBA color array on the LEDs.
func displayLEDArray(leds *ws2812.Device, ledlist [6]color.RGBA) error {
	err := leds.WriteColors(ledlist[:])
	return err
}

// checkInputCols checks the input columns for a button press.
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
	// Sleep to give the user time to read the error
	time.Sleep(2 * time.Second)
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
