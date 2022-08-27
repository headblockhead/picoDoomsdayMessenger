package picodoomsdaymessenger

import (
	"errors"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type State string

const (
	StateMenu      State = "menu"
	StatePeople    State = "people"
	StateMessanges State = "messanges"
	StateSettings  State = "settings"
	StateShutdown  State = "shutdown"
)

var (
	MenuMenuTitle          = "Main Menu"
	MenuMenuItemMessages   = "Messages"
	MenuMenuItemPeople     = "People"
	MenuMenuItemSettings   = "Settings"
	MenuMenuItemShutdown   = "Shutdown"
	GlobalMenuItemGoBack   = "Go Back"
	MenuMenuItems          = []string{MenuMenuItemMessages, MenuMenuItemPeople, MenuMenuItemSettings, MenuMenuItemShutdown}
	MessagesMenuTitle      = "Messages"
	MessagesMenuItems      = []string{GlobalMenuItemGoBack}
	PeopleMenuTitle        = "People"
	PeopleMenuItems        = []string{GlobalMenuItemGoBack}
	SettingsMenuTitle      = "Settings"
	SettingsMenuItems      = []string{GlobalMenuItemGoBack}
	ShutdownMenuTitle      = "Are you sure?"
	ShutdownMenuItemReally = "Yes, really."
	ShutdownMenuItems      = []string{GlobalMenuItemGoBack, ShutdownMenuItemReally}
)

type Device struct {
	State           State
	StateHistory    []State
	MenuTitle       string
	MenuItems       []string
	CurrentMenuItem int
}

type InputEvent string

const (
	InputEventLeft  InputEvent = "left"
	InputEventRight InputEvent = "right"
	InputEventFire  InputEvent = "fire"
)

func NewDevice() (device *Device) {
	return &Device{StateMenu, []State{}, MenuMenuTitle, MenuMenuItems, 0}
}

func (m *Device) ChangeState(newState State) {
	m.StateHistory = append(m.StateHistory, newState)
	m.State = newState
}

func (m *Device) ChangeStateWithoutHistory(newState State) {
	m.State = newState
}

func (m *Device) GoBackState() (err error) {
	if len(m.StateHistory)-1 < 0 {
		return errors.New("already at root state")
	}
	m.ChangeStateWithoutHistory(m.StateHistory[len(m.StateHistory)-1])
	return nil
}

func (m *Device) ProcessInputEvent(event InputEvent) (err error) {
	switch m.State {
	case StateMenu:
		{
			m.checkMenuScroll(event)
			if event == InputEventFire {
				switch m.MenuItems[m.CurrentMenuItem] {
				case MenuMenuItemMessages:
					{
						m.ChangeState(StateMessanges)
					}
				case MenuMenuItemPeople:
					{
						m.ChangeState(StatePeople)
					}
				case MenuMenuItemSettings:
					{
						m.ChangeState(StateSettings)
					}
				case MenuMenuItemShutdown:
					{
						m.ChangeState(StateShutdown)
					}
				}
			}
		}
	case StateMessanges:
		{
			m.checkMenuScroll(event)
			if event == InputEventFire {
				if m.MenuItems[m.CurrentMenuItem] == GlobalMenuItemGoBack {
					err = m.GoBackState()
					if err != nil {
						return err
					}
				}
			}
		}
	case StatePeople:
		{
			m.checkMenuScroll(event)
			if event == InputEventFire {
				if m.MenuItems[m.CurrentMenuItem] == GlobalMenuItemGoBack {
					err = m.GoBackState()
					if err != nil {
						return err
					}
				}
			}
		}
	case StateSettings:
		{
			m.checkMenuScroll(event)
			if event == InputEventFire {
				if m.MenuItems[m.CurrentMenuItem] == GlobalMenuItemGoBack {
					err = m.GoBackState()
					if err != nil {
						return err
					}
				}
			}
		}
	case StateShutdown:
		{
			m.checkMenuScroll(event)
			if event == InputEventFire {
				if m.MenuItems[m.CurrentMenuItem] == GlobalMenuItemGoBack {
					err = m.GoBackState()
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (m *Device) checkMenuScroll(event InputEvent) {
	if event == InputEventLeft {
		if m.CurrentMenuItem > 0 {
			m.CurrentMenuItem--
		}
	}
	if event == InputEventRight {
		if m.CurrentMenuItem < len(m.MenuItems) {
			m.CurrentMenuItem++
		}
	}
}

func GetFrame(dimensions image.Rectangle, device *Device) (frame image.Image, err error) {
	img := image.NewRGBA(dimensions)
	drawText(img, 0, 13, device.MenuTitle)
	drawHLine(img, 0, 15, dimensions.Dx())
	for i := 0; i < len(device.MenuItems); i++ {
		drawText(img, 0, 26+(13*(i)), device.MenuItems[i])
	}
	drawCursor(img, dimensions.Dx()-4, 6+(13*(device.CurrentMenuItem+1)))
	return img, nil
}

func GetErrorFrame(dimensions image.Rectangle, device *Device, inputErr string) (frame image.Image, err error) {
	img := image.NewRGBA(dimensions)
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

func drawHLine(img *image.RGBA, x1 int, y int, x2 int) {
	col := color.RGBA{255, 255, 255, 255}
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

func drawVLine(img *image.RGBA, y1 int, x int, y2 int) {
	col := color.RGBA{255, 255, 255, 255}
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}
