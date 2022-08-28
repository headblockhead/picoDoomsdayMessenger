package main

import (
	"fmt"
	"image"
	"reflect"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	picodoomsdaymessenger "github.com/headblockhead/picoDoomsdayMessenger"
	"github.com/nfnt/resize"
	"golang.org/x/image/colornames"
)

var currentFrame image.Image

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "DoomsDayMessenger",
		Bounds: pixel.R(0, 0, 512, 256),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Black)

	// Create a new Machine
	device := picodoomsdaymessenger.NewDevice()
	// Record the display size
	displayx, displayy := 128, 64
	// Set the old machine state and old menu item to something that is not the starting value.
	oldDeviceState := picodoomsdaymessenger.StateUnknown
	oldDeviceHighlightedItem := picodoomsdaymessenger.GlobalMenuItemUnknown

	// Panic recovery
	defer func() {
		if err := recover(); err != nil {
			// If there is a panic, try to print details to the screen.
			// The handleError() function cannot be used here as it requires an error.
			// err in this case is not an error but an interface.
			// So we use fmt.Sprintf("%v", err) to write details to the screen.
			frame, newErr := picodoomsdaymessenger.GetErrorFrame(image.Rect(0, 0, int(displayx), int(displayy)), device, fmt.Sprintf("%v", err))
			if newErr != nil {
				return
			}
			newErr = displayImage(win, frame)
			if newErr != nil {
				return
			}
		}
	}()

	for !win.Closed() {
		// Update the display if the state changes
		if !reflect.DeepEqual(oldDeviceState, device.State) || !reflect.DeepEqual(oldDeviceHighlightedItem, device.State.HighlightedItem) {
			oldDeviceState = device.State
			oldDeviceHighlightedItem = device.State.HighlightedItem
			frame, err := picodoomsdaymessenger.GetFrame(image.Rect(0, 0, int(displayx), int(displayy)), device)
			if err != nil {
				handleError(win, device, err)
				return
			}
			err = displayImage(win, frame)
			if err != nil {
				handleError(win, device, err)
				return
			}
		}
		// Send a message to the device if an input is triggered.

		// Rotary encoder checking
		// If the rotary encoder is turned
		if win.JustPressed(pixelgl.KeyUp) {
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventFire)
			if err != nil {
				handleError(win, device, err)
				return
			}
			// Wait before checking again to make sure we don't take in multiple clicks.
			time.Sleep(time.Millisecond * 100)
		}
		if win.JustPressed(pixelgl.KeyLeft) {
			// Anti-Clockwise
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventLeft)
			if err != nil {
				// If there is a error, try to print details to the screen.
				handleError(win, device, err)
				return
			}
		} else if win.JustPressed(pixelgl.KeyRight) {
			// Clockwise
			err := device.ProcessInputEvent(picodoomsdaymessenger.InputEventRight)
			if err != nil {
				// If there is a error, try to print details to the screen.
				handleError(win, device, err)
				return
			}
		}
		win.Update()
		time.Sleep(time.Millisecond * 1)
	}
}

// displayImage takes in an image and writes it to the screen.
func displayImage(win *pixelgl.Window, img image.Image) (err error) {
	pixelArray := []uint8{}
	img = resize.Resize(uint(win.Bounds().Max.X), uint(win.Bounds().Max.Y), img, resize.NearestNeighbor)
	// Put the image into the buffer.
	for y := img.Bounds().Dy(); y > 0; y-- {
		for x := 0; x < img.Bounds().Dx(); x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if uint8(r) == 255 && uint8(g) == 255 && uint8(b) == 255 && uint8(a) == 255 {
				if y < 68 {
					pixelArray = append(pixelArray, uint8(255))
					pixelArray = append(pixelArray, uint8(170))
					pixelArray = append(pixelArray, uint8(0))
					pixelArray = append(pixelArray, uint8(255))
				} else {
					pixelArray = append(pixelArray, uint8(20))
					pixelArray = append(pixelArray, uint8(240))
					pixelArray = append(pixelArray, uint8(255))
					pixelArray = append(pixelArray, uint8(255))
				}
			} else {
				pixelArray = append(pixelArray, uint8(r))
				pixelArray = append(pixelArray, uint8(g))
				pixelArray = append(pixelArray, uint8(b))
				pixelArray = append(pixelArray, uint8(a))
			}
		}
	}
	win.Canvas().SetPixels(pixelArray)
	return nil
}

// handleError takes in an error and communicates it to the user.
func handleError(win *pixelgl.Window, device *picodoomsdaymessenger.Device, inputerr error) {
	// Try to get details to print to the screen
	displayx, displayy := 128, 64
	frame, newErr := picodoomsdaymessenger.GetErrorFrame(image.Rect(0, 0, int(displayx), int(displayy)), device, inputerr.Error())
	if newErr != nil {
		// If we can't do that, resort to signaling with the LED
		return
	}
	// Try to print the details to the screen
	newErr = displayImage(win, frame)
	if newErr != nil {
		// If we can't do that either, resort to signaling with the LED
		return
	}
}

func main() {
	pixelgl.Run(run)
}
