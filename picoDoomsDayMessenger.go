package picodoomsdaymessenger

import (
	"errors"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Basic structure.
type Device struct {
	State        *State
	StateHistory []*State
}
type State struct {
	Title           string
	Content         []MenuItem
	HighlightedItem MenuItem
	LoadAction      func(d *Device) (err error)
}
type MenuItem struct {
	Name   string
	Action func(d *Device) (err error)
	Index  int
}

type InputEvent string

const (
	InputEventUp   InputEvent = "left"
	InputEventDown InputEvent = "right"
	InputEventFire InputEvent = "fire"
)

// Define MenuItems
var (
	GlobalMenuItemDefault MenuItem = MenuItem{
		Name: "DefaultMenuItem",
		Action: func(d *Device) (err error) {
			return errors.New("If you see this error, this is not intended. Please contact the dev.")
		},
		Index: 0,
	}
	GlobalMenuItemGoBack MenuItem = MenuItem{
		Name: "Go Back",
		Action: func(d *Device) (err error) {
			err = d.GoBackState()
			if err != nil {
				return err
			}
			return nil
		},
		Index: 0,
	}
	MainMenuItemMessages MenuItem = MenuItem{
		Name: "Messages",
		Action: func(d *Device) (err error) {
			d.ChangeStateWithHistory(&StateMessagesMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index: 0,
	}
	MainMenuItemPeople MenuItem = MenuItem{
		Name: "People",
		Action: func(d *Device) (err error) {
			d.ChangeStateWithHistory(&StatePeopleMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index: 1,
	}
	MainMenuItemSettings MenuItem = MenuItem{
		Name: "Settings",
		Action: func(d *Device) (err error) {
			d.ChangeStateWithHistory(&StateSettingsMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index: 2,
	}
	MainMenuItemSleep MenuItem = MenuItem{
		Name: "Sleep",
		Action: func(d *Device) (err error) {
			return errors.New("Something has gone terribly wrong here!")
		},
		Index: 3,
	}
)

// Define States
var (
	StateDefault = State{
		Title:           "DefaultState",
		Content:         []MenuItem{GlobalMenuItemDefault},
		HighlightedItem: GlobalMenuItemDefault,
	}
	StateMainMenu = State{
		Title:           "Main Menu",
		Content:         []MenuItem{MainMenuItemMessages, MainMenuItemPeople, MainMenuItemSettings, MainMenuItemSleep},
		HighlightedItem: MainMenuItemMessages,
	}
	StateMessagesMenu = State{
		Title:           "Messages",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: GlobalMenuItemGoBack,
	}
	StatePeopleMenu = State{
		Title:           "People",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: GlobalMenuItemGoBack,
	}
	StateSettingsMenu = State{
		Title:           "Settings",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: GlobalMenuItemGoBack,
	}
)

// NewDevice returns a new Device with default parameters.
func NewDevice() (d *Device) {
	newDevice := &Device{&StateMainMenu, []*State{&StateMainMenu}}
	return newDevice
}

// ChangeStateWithHistory will take in a State and update the Device while adding the State to the StateHistory.
func (d *Device) ChangeStateWithHistory(newState *State) (err error) {
	d.StateHistory = append(d.StateHistory, newState)
	err = d.ChangeStateWithoutHistory(newState)
	if err != nil {
		return err
	}
	return nil
}

// ChangeStateWithoutHistory will take in a State and update the Device.
func (d *Device) ChangeStateWithoutHistory(newState *State) (err error) {
	d.State = newState
	if d.State.LoadAction != nil {
		err = d.State.LoadAction(d)
		if err != nil {
			return err
		}
	}
	return nil
}

// GoBackState will use the StateHistory to return to the upwards state in the tree.
func (d *Device) GoBackState() (err error) {
	if len(d.StateHistory) < 1 {
		return errors.New("already at root state")
	}
	err = d.ChangeStateWithoutHistory(d.StateHistory[len(d.StateHistory)-2])
	if err != nil {
		return err
	}
	d.StateHistory = d.StateHistory[0 : len(d.StateHistory)-1]
	return nil
}

// ProcessInputEvent will take in an InputEvent and run appropriate actions based on the event.
func (d *Device) ProcessInputEvent(inputEvent InputEvent) (err error) {
	switch inputEvent {
	case InputEventUp:
		{
			if d.State.HighlightedItem.Index > 0 {
				d.State.HighlightedItem = d.State.Content[d.State.HighlightedItem.Index-1]
			}
		}
	case InputEventDown:
		{
			if d.State.HighlightedItem.Index < len(d.State.Content)-1 {
				d.State.HighlightedItem = d.State.Content[d.State.HighlightedItem.Index+1]
			}
		}
	case InputEventFire:
		{
			err = d.State.HighlightedItem.Action(d)
			return err
		}
	}
	return nil
}

// GetFrame will take in a Device and return an image based on the state.
func GetFrame(dimensions image.Rectangle, d *Device) (frame image.Image, err error) {
	img := image.NewRGBA(dimensions)
	drawText(img, 0, 13, d.State.Title)
	drawHLine(img, 0, 15, dimensions.Dx())
	for i := 0; i < len(d.State.Content); i++ {
		drawText(img, 0, 26+(13*(i)), d.State.Content[i].Name)
	}
	drawCursor(img, dimensions.Dx()-4, 6+(13*(d.State.HighlightedItem.Index+1)))
	return img, nil
}

// GetErrorFrame will take in a string version of an error and return an image with that error in.
func GetErrorFrame(dimensions image.Rectangle, d *Device, inputErr string) (frame image.Image, err error) {
	img := image.NewRGBA(dimensions)
	inputErr = "FATAL ERR: " + inputErr
	if len(inputErr) < 18 {
		drawText(img, 0, 13, inputErr)
	} else if len(inputErr) > 18 && len(inputErr) < 36 {
		drawText(img, 0, 13, inputErr[:18])
		drawText(img, 0, 26, inputErr[18:])
	} else if len(inputErr) > 36 && len(inputErr) < 54 {
		drawText(img, 0, 13, inputErr[:18])
		drawText(img, 0, 26, inputErr[18:36])
		drawText(img, 0, 39, inputErr[36:])
	} else {
		drawText(img, 0, 13, inputErr[:18])
		drawText(img, 0, 26, inputErr[18:36])
		drawText(img, 0, 39, inputErr[36:54])
		drawText(img, 0, 52, inputErr[54:])
	}
	return img, nil
}

// drawText will write text in a 7x13 pixel font at a location.
func drawText(img *image.RGBA, x, y int, text string) {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

// drawCursor will draw a small arrow. It is drawn based on the X and Y coordinates being at the top-left of the arrow.
func drawCursor(img *image.RGBA, x int, y int) {
	col := color.RGBA{255, 255, 255, 255}
	img.Set(x+0, y+0, col)
	img.Set(x+1, y+1, col)
	img.Set(x+2, y+2, col)
	img.Set(x+3, y+3, col)
	img.Set(x+2, y+4, col)
	img.Set(x+1, y+5, col)
	img.Set(x+0, y+6, col)
}

// drawHLine draws a horizontal line from one X location to another. x2 has to be greater than x1.
func drawHLine(img *image.RGBA, x1 int, y int, x2 int) {
	col := color.RGBA{255, 255, 255, 255}
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// drawVLine draws a verticle line from one Y location to another. y2 has to be greater than y1.
func drawVLine(img *image.RGBA, y1 int, x int, y2 int) {
	col := color.RGBA{255, 255, 255, 255}
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}
