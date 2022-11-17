package picodoomsdaymessenger

import (
	"errors"
	"image"
	"image/color"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Device is the main structure that holds all the information about the device. It has a State, a StateHistory, and an LEDAnimation.
type Device struct {
	State               *State
	StateHistory        []*State
	LEDAnimation        *LEDAnimation
	Conversations       []*Conversation
	CurrentConversation *Conversation
	SelfIdentity        *Person
}

// Conversation is a conversation with a person. It contains a list of Messages and a Person that the conversation is with.
type Conversation struct {
	Messages           []Message
	HighlightedMessage *Message
	Name               string
}

// Person is a representation of another device. A Person has a name and a unique identifier
type Person struct {
	Name string
	ID   uint32
}

// Message is a message sent inside a Conversation. It contains the time it was sent, the time it was recieved and the content of the message.
type Message struct {
	TimeSent     time.Time
	TimeRecieved time.Time
	Text         string
	Index        int
	Person       Person
}

// State is the current state of the device. It contains all the information about what is currently being displayed.
type State struct {
	Title           string
	Content         []MenuItem
	HighlightedItem *MenuItem
	LoadAction      func(d *Device) (err error)
}

// MenuItem is a structure that holds data that can be displayed on the screen. It contains a title and an action that is run when the item is selected.
type MenuItem struct {
	Text          string
	Action        func(d *Device) (err error)
	Index         int
	GetCursorData func(d *Device) (data any, err error)
	CursorIcon    CursorIcon
}

// CursorIcon is a function that draws a cursor icon based on the data at a location.
type CursorIcon func(img *image.RGBA, x int, y int, data any) (err error)

// LEDAnimation is a structure that holds information about an LED animation.
type LEDAnimation struct {
	FrameDuration time.Duration
	CurrentFrame  int
	Frames        [][6]color.RGBA
}

// Define default People
var PersonDefault = Person{"You", 0}

// Define Cursors
var (
	// CursorIconRightArrow is a cursor that is a right arrow. It does not need any data.
	CursorIconRightArrow = func(img *image.RGBA, x int, y int, data any) (err error) {
		col := color.RGBA{255, 255, 255, 255}
		img.Set(x+0, y+0, col)
		img.Set(x+1, y+1, col)
		img.Set(x+2, y+2, col)
		img.Set(x+3, y+3, col)
		img.Set(x+2, y+4, col)
		img.Set(x+1, y+5, col)
		img.Set(x+0, y+6, col)
		return nil
	}
	// CursorIconLeftArrow is a cursor that is a left arrow. It does not need any data.
	CursorIconLeftArrow = func(img *image.RGBA, x int, y int, data any) (err error) {
		col := color.RGBA{255, 255, 255, 255}
		img.Set(x+6, y+0, col)
		img.Set(x+5, y+1, col)
		img.Set(x+4, y+2, col)
		img.Set(x+3, y+3, col)
		img.Set(x+4, y+4, col)
		img.Set(x+5, y+5, col)
		img.Set(x+6, y+6, col)
		return nil
	}
	// CursorIconBox is a cursor that is a box. It takes in a bool as data. If the bool is true, the box will be filled in. If the bool is false, the box will be empty.
	CursorIconBox = func(img *image.RGBA, x int, y int, data any) (err error) {
		isChecked, ok := data.(bool)
		if !ok {
			return errors.New("data is not a bool")
		}
		for i := 0; i < 7; i++ {
			for j := 0; j < 7; j++ {
				img.Set(x+i, y+j, color.RGBA{255, 255, 255, 255})
			}
		}
		if !isChecked {
			for i := 1; i < 6; i++ {
				for j := 1; j < 6; j++ {
					img.Set(x+i, y+j, color.RGBA{0, 0, 0, 255})
				}
			}
		}
		return nil
	}
)

// Define MenuItems
var (

	// Default Menu Items

	// MenuItemDefault is a MenuItem that does nothing. It is used as a placeholder for the default State of the device.
	MenuItemDefault MenuItem = MenuItem{
		Text: "DefaultMenuItem",
		Action: func(d *Device) (err error) {
			return errors.New("default menu item action")
		},
		Index:      0,
		CursorIcon: CursorIconRightArrow,
	}

	// Global Menu Items

	// GlobalMenuItemGoBack is a MenuItem that goes back to the previous state in the StateHistory.
	GlobalMenuItemGoBack MenuItem = MenuItem{
		Text: "Go Back",
		Action: func(d *Device) (err error) {
			err = d.GoBackState()
			if err != nil {
				return err
			}
			return nil
		},
		Index:      0,
		CursorIcon: CursorIconLeftArrow,
	}

	// Main Menu Items

	// MainMenuItemMessages is a MenuItem that goes to the Messages menu.
	MainMenuItemMessages MenuItem = MenuItem{
		Text: "Messages",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StateMessagesMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index:      0,
		CursorIcon: CursorIconRightArrow,
	}

	// MainMenuItemPeople is a MenuItem that goes to the People menu.
	MainMenuItemPeople MenuItem = MenuItem{
		Text: "People",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StatePeopleMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index:      1,
		CursorIcon: CursorIconRightArrow,
	}

	// MainMenuItemGames is a MenuItem that goes to the Games menu.
	MainMenuItemGames MenuItem = MenuItem{
		Text: "Games",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StateGamesMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index:      2,
		CursorIcon: CursorIconRightArrow,
	}

	// MainMenuItemDemos is a MenuItem that goes to the Demos menu.
	MainMenuItemDemos MenuItem = MenuItem{
		Text: "Demo",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StateDemosMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index:      3,
		CursorIcon: CursorIconRightArrow,
	}

	// MainMenuItemTools is a MenuItem that goes to the Tools menu.
	MainMenuItemTools MenuItem = MenuItem{
		Text: "Tools",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StateToolsMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index:      4,
		CursorIcon: CursorIconRightArrow,
	}

	// MainMenuItemSettings is a MenuItem that goes to the Settings menu.
	MainMenuItemSettings MenuItem = MenuItem{
		Text: "Settings",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StateSettingsMenu)
			if err != nil {
				return err
			}
			return nil
		},
		Index:      5,
		CursorIcon: CursorIconRightArrow,
	}

	// Games Menu Items

	// Demos Menu Items

	// DemoMenuItemRGB is a MenuItem that toggles a demo of the RGB LEDs.
	DemoMenuItemRGB MenuItem = MenuItem{
		Text: "RGB Demo",
		Action: func(d *Device) (err error) {
			if d.LEDAnimation != &LEDAnimationDemo {
				err = d.ChangeLEDAnimationWithoutContinue(&LEDAnimationDemo)
				if err != nil {
					return err
				}
			} else {
				err = d.ChangeLEDAnimationWithoutContinue(&LEDAnimationDefault)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Index: 1,
		GetCursorData: func(d *Device) (data any, err error) {
			return d.LEDAnimation == &LEDAnimationDemo, nil
		},
		CursorIcon: CursorIconBox,
	}

	// Tools Menu Items

	// ToolsMenuItemSOS is a MenuItem that toggles a SOS message shown in morse code through the RGB LEDs.
	ToolsMenuItemSOS MenuItem = MenuItem{
		Text: "SOS Mode",
		Action: func(d *Device) (err error) {
			if d.LEDAnimation != &LEDAnimationSOS {
				err = d.ChangeLEDAnimationWithoutContinue(&LEDAnimationSOS)
				if err != nil {
					return err
				}
			} else {
				err = d.ChangeLEDAnimationWithoutContinue(&LEDAnimationDefault)
				if err != nil {
					return err
				}
			}
			return nil
		},
		Index: 1,
		GetCursorData: func(d *Device) (data any, err error) {
			return d.LEDAnimation == &LEDAnimationSOS, nil
		},
		CursorIcon: CursorIconBox,
	}
)

// Define States
var (
	// StateDefault is a State that does nothing. It is used as a placeholder for the default State of the Device.
	StateDefault = State{
		Title:           "DefaultState",
		Content:         []MenuItem{MenuItemDefault},
		HighlightedItem: &MenuItemDefault,
	}
	// StateConversationReader is a special State that is used when reading a Conversation.
	StateConversationReader = State{
		Title:   "",
		Content: []MenuItem{},
	}
	// StateMainMenu is a State that shows the main menu.
	StateMainMenu = State{
		Title:           "Main Menu",
		Content:         []MenuItem{MainMenuItemMessages, MainMenuItemPeople, MainMenuItemGames, MainMenuItemDemos, MainMenuItemTools, MainMenuItemSettings},
		HighlightedItem: &MainMenuItemMessages,
	}
	// StateMessagesMenu is a State that shows the messages menu.
	StateMessagesMenu = State{
		Title:           "Messages",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: &GlobalMenuItemGoBack,
	}
	// StateMessagesMenuOld is a copy of StateMessagesMenu that can be used as a starting point to reset StateMessagesMenu.
	StateMessagesMenuOld = StateMessagesMenu
	// StatePeopleMenu is a State that shows the people menu.
	StatePeopleMenu = State{
		Title:           "People",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: &GlobalMenuItemGoBack,
	}
	// StateGamesMenu is a State that shows the games menu.
	StateGamesMenu = State{
		Title:           "Games",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: &GlobalMenuItemGoBack,
	}
	// StateDemosMenu is a State that shows the demos menu.
	StateDemosMenu = State{
		Title:           "Demos",
		Content:         []MenuItem{GlobalMenuItemGoBack, DemoMenuItemRGB},
		HighlightedItem: &GlobalMenuItemGoBack,
	}
	// StateToolsMenu is a State that shows the tools menu.
	StateToolsMenu = State{
		Title:           "Tools",
		Content:         []MenuItem{GlobalMenuItemGoBack, ToolsMenuItemSOS},
		HighlightedItem: &GlobalMenuItemGoBack,
	}
	// StateSettingsMenu is a State that shows the settings menu.
	StateSettingsMenu = State{
		Title:           "Settings",
		Content:         []MenuItem{GlobalMenuItemGoBack},
		HighlightedItem: &GlobalMenuItemGoBack,
	}
)

// Define LED animations. They are made of multiple frames of 6 colors.
var (
	// LEDAnimationDefault is the default LED animation. It is used when no other animation is active and is simply black.
	LEDAnimationDefault = LEDAnimation{
		FrameDuration: 5 * time.Millisecond,
		CurrentFrame:  0,
		Frames: [][6]color.RGBA{
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
		},
	}
	// LEDAnimationSOS is an LED animation that shows the SOS message in morse code.
	LEDAnimationSOS = LEDAnimation{
		FrameDuration: 200 * time.Millisecond,
		CurrentFrame:  0,
		Frames: [][6]color.RGBA{
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 255}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
		},
	}
	// LEDAnimationDemo is an LED animation that shows off the capabilities of the LED animation system.
	LEDAnimationDemo = LEDAnimation{
		FrameDuration: 1 * time.Millisecond,
		CurrentFrame:  0,
		Frames: [][6]color.RGBA{
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 255, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 255, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 255, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{0, 0, 255, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{255, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 255, 0}},

			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},

			{color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}},
			{color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}},
			{color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}},
			{color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}},
			{color.RGBA{255, 255, 255, 0}, color.RGBA{255, 255, 255, 0}, color.RGBA{255, 255, 255, 0}, color.RGBA{255, 255, 255, 0}, color.RGBA{255, 255, 255, 0}, color.RGBA{255, 255, 255, 0}},
			{color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}, color.RGBA{200, 200, 200, 0}},
			{color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}, color.RGBA{150, 150, 150, 0}},
			{color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}, color.RGBA{100, 100, 100, 0}},
			{color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}, color.RGBA{50, 50, 50, 0}},

			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},

			{color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}},
			{color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}},
			{color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}},
			{color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}},
			{color.RGBA{255, 0, 0, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{255, 0, 0, 0}, color.RGBA{255, 0, 0, 0}},
			{color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}, color.RGBA{200, 0, 0, 0}},
			{color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}, color.RGBA{150, 0, 0, 0}},
			{color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}, color.RGBA{100, 0, 0, 0}},
			{color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}, color.RGBA{50, 0, 0, 0}},

			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},

			{color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}},
			{color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}},
			{color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}},
			{color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}},
			{color.RGBA{0, 255, 0, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 255, 0, 0}, color.RGBA{0, 255, 0, 0}},
			{color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}, color.RGBA{0, 200, 0, 0}},
			{color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}, color.RGBA{0, 150, 0, 0}},
			{color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}, color.RGBA{0, 100, 0, 0}},
			{color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}, color.RGBA{0, 050, 0, 0}},

			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},

			{color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}},
			{color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}},
			{color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}},
			{color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}},
			{color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 255, 0}, color.RGBA{0, 0, 255, 0}},
			{color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}, color.RGBA{0, 0, 200, 0}},
			{color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}, color.RGBA{0, 0, 150, 0}},
			{color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}, color.RGBA{0, 0, 100, 0}},
			{color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}, color.RGBA{0, 0, 050, 0}},

			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
			{color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}, color.RGBA{0, 0, 0, 0}},
		}}
)

// NewDevice returns a new Device with default parameters.
func NewDevice() (d *Device) {
	return &Device{&StateMainMenu, []*State{&StateMainMenu}, &LEDAnimationDefault, []*Conversation{}, &Conversation{}, &PersonDefault}
}

// NewConversation creates a blank new Conversation and adds it to the Device. It also returns a pointer to that Conversation.
func (d *Device) NewConversation() (c *Conversation) {
	newConversation := &Conversation{}
	d.Conversations = append(d.Conversations, newConversation)
	return newConversation
}

func (d *Device) UpdateMessagesMenu() {
	StateMessagesMenu = StateMessagesMenuOld
	for i := 0; i < len(d.Conversations); i++ {
		// Define a seperate variable to seperate the increasing i from the functions defined here.
		j := i
		StateMessagesMenu.Content = append(StateMessagesMenu.Content, MenuItem{
			Text: d.Conversations[j].Name,
			Action: func(d *Device) (err error) {
				d.CurrentConversation = d.Conversations[j]
				err = d.ChangeStateWithHistory(&StateConversationReader)
				return err
			},
			Index:      j + 1,
			CursorIcon: CursorIconRightArrow,
		})
	}
}

// ChangeLEDAnimationWithoutContinue changes the current LED animation of the device without continuing from the last time it was played.
func (d *Device) ChangeLEDAnimationWithoutContinue(newAnimation *LEDAnimation) (err error) {
	d.LEDAnimation = newAnimation
	d.LEDAnimation.CurrentFrame = 0
	return nil
}

// ChangeLEDAnimation changes the current LED animation of the device and continues from the last time it was played.
func (d *Device) ChangeLEDAnimationWithContinue(newAnimation *LEDAnimation) (err error) {
	d.LEDAnimation = newAnimation
	return nil
}

// ChangeStateWithHistory will take in a State and update the Device while adding the State to the StateHistory.
func (d *Device) ChangeStateWithHistory(newState *State) (err error) {
	d.StateHistory = append(d.StateHistory, newState)
	err = d.ChangeStateWithoutHistory(newState)
	return err
}

// ChangeStateWithoutHistory will take in a State and update the Device.
func (d *Device) ChangeStateWithoutHistory(newState *State) (err error) {
	d.State = newState
	if d.State.LoadAction != nil {
		err = d.State.LoadAction(d)
	}
	return err
}

// GoBackState will use the StateHistory to return to the upwards state in the tree.
func (d *Device) GoBackState() (err error) {
	if len(d.StateHistory) <= 1 {
		return errors.New("already at root state")
	}
	err = d.ChangeStateWithoutHistory(d.StateHistory[len(d.StateHistory)-2])
	d.StateHistory = d.StateHistory[0 : len(d.StateHistory)-1]
	return err
}

// InputEvent is a string that represents a button press.
type InputEvent string

const (
	InputEventUp           InputEvent = "up"
	InputEventDown         InputEvent = "down"
	InputEventLeft         InputEvent = "left"
	InputEventRight        InputEvent = "right"
	InputEventAccept       InputEvent = "accept"
	InputEventFunction1    InputEvent = "function1"
	InputEventFunction2    InputEvent = "function2"
	InputEventFunction3    InputEvent = "function3"
	InputEventFunction4    InputEvent = "function4"
	InputEventOpenSettings InputEvent = "openSettings"
	InputEventOpenPeople   InputEvent = "openPeople"
	InputEventOpenMessages InputEvent = "openMessages"
	InputEventOpenMainMenu InputEvent = "openMainMenu"
	InputEventNumber1      InputEvent = "number1"
	InputEventNumber2      InputEvent = "number2"
	InputEventNumber3      InputEvent = "number3"
	InputEventNumber4      InputEvent = "number4"
	InputEventNumber5      InputEvent = "number5"
	InputEventNumber6      InputEvent = "number6"
	InputEventNumber7      InputEvent = "number7"
	InputEventNumber8      InputEvent = "number8"
	InputEventNumber9      InputEvent = "number9"
	InputEventNumber0      InputEvent = "number0"
	InputEventStar         InputEvent = "star"
	InputEventPound        InputEvent = "pound"
)

// ProcessInputEvent will take in an InputEvent and run appropriate actions based on the event.
func (d *Device) ProcessInputEvent(inputEvent InputEvent) (err error) {
	switch inputEvent {
	case InputEventUp:
		{
			if d.State != &StateConversationReader {
				if d.State.HighlightedItem.Index <= 0 {
					d.State.HighlightedItem = &d.State.Content[len(d.State.Content)-1]
				} else {
					d.State.HighlightedItem = &d.State.Content[d.State.HighlightedItem.Index-1]
				}
			} else {
				if d.CurrentConversation.HighlightedMessage.Index <= 0 {
					d.CurrentConversation.HighlightedMessage = &d.CurrentConversation.Messages[len(d.CurrentConversation.Messages)-1]
				} else {
					d.CurrentConversation.HighlightedMessage = &d.CurrentConversation.Messages[d.CurrentConversation.HighlightedMessage.Index-1]
				}
			}
		}
	case InputEventDown:
		{
			if d.State != &StateConversationReader {
				if d.State.HighlightedItem.Index >= len(d.State.Content)-1 {
					d.State.HighlightedItem = &d.State.Content[0]
				} else {
					d.State.HighlightedItem = &d.State.Content[d.State.HighlightedItem.Index+1]
				}
			} else {
				if d.CurrentConversation.HighlightedMessage.Index >= len(d.CurrentConversation.Messages)-1 {
					d.CurrentConversation.HighlightedMessage = &d.CurrentConversation.Messages[0]
				} else {
					d.CurrentConversation.HighlightedMessage = &d.CurrentConversation.Messages[d.CurrentConversation.HighlightedMessage.Index+1]
				}
			}
		}
	case InputEventAccept:
		{
			if d.State != &StateConversationReader {
				err = d.State.HighlightedItem.Action(d)
				return err
			} else {
				return errors.New("cannot accept in conversation reader")
			}
		}
	case InputEventOpenSettings:
		{
			err = d.ChangeStateWithHistory(&StateSettingsMenu)
			return err
		}
	case InputEventOpenPeople:
		{
			err = d.ChangeStateWithHistory(&StatePeopleMenu)
			return err
		}
	case InputEventOpenMessages:
		{
			err = d.ChangeStateWithHistory(&StateMessagesMenu)
			return err
		}
	case InputEventOpenMainMenu:
		{
			err = d.ChangeStateWithHistory(&StateMainMenu)
			return err
		}
	}
	return nil
}

// GetFrame will take in a Device and return an image based on the state.
func GetFrame(dimensions image.Rectangle, d *Device) (frame image.Image, err error) {
	img := image.NewRGBA(dimensions)

	if d.State != &StateConversationReader {
		// Draw the content with the currently highlighted item in the middle of the screen and the other items above and below it.
		for i := 0; i < len(d.State.Content); i++ {
			if d.State.Content[i].Index == d.State.HighlightedItem.Index {
				drawText(img, 0, 43, d.State.Content[i].Text)
			} else if d.State.Content[i].Index < d.State.HighlightedItem.Index {
				drawText(img, 0, 43-(d.State.HighlightedItem.Index-d.State.Content[i].Index)*12, d.State.Content[i].Text)
			} else if d.State.Content[i].Index > d.State.HighlightedItem.Index {
				drawText(img, 0, 43+(d.State.Content[i].Index-d.State.HighlightedItem.Index)*12, d.State.Content[i].Text)
			}
		}

		// Draw the title.
		drawBlackFilledBox(img, 0, 0, dimensions.Dx(), 16)
		drawText(img, 0, 13, d.State.Title)
		drawHLine(img, 0, 15, dimensions.Dx())

		// Draw the cursor. If the cursor is a checkbox, check if the checkbox is checked or not.
		var cursorData any
		if d.State.HighlightedItem.GetCursorData != nil {
			cursorData, err = d.State.HighlightedItem.GetCursorData(d)
			if err != nil {
				return nil, err
			}
		}
		err = d.State.HighlightedItem.CursorIcon(img, dimensions.Dx()-7, 36, cursorData)
		if err != nil {
			return nil, err
		}
	} else {
		drawBlackFilledBox(img, 0, 0, (dimensions.Dx()*75)/100, 16)
		drawText(img, 0, 13, d.CurrentConversation.Name)
		drawHLine(img, 0, 15, dimensions.Dx()*75)

		// Draw the conversation with the most recent message at the bottom of the screen.
		for i := 0; i < len(d.CurrentConversation.Messages); i++ {
			if d.CurrentConversation.Messages[i].Index == d.CurrentConversation.HighlightedMessage.Index {
				if d.CurrentConversation.Messages[i].Person != *d.SelfIdentity {
					drawText(img, 0, 43, d.CurrentConversation.Messages[i].Text)
				} else {
					drawText(img, dimensions.Dx()-(len(d.CurrentConversation.Messages[i].Text)*7), 43, d.CurrentConversation.Messages[i].Text)
				}
			} else if d.CurrentConversation.Messages[i].Index < d.CurrentConversation.HighlightedMessage.Index {
				if d.CurrentConversation.Messages[i].Person != *d.SelfIdentity {
					drawText(img, 0, 43-(d.CurrentConversation.HighlightedMessage.Index-d.CurrentConversation.Messages[i].Index)*12, d.CurrentConversation.Messages[i].Text)
				} else {
					drawText(img, dimensions.Dx()-(len(d.CurrentConversation.Messages[i].Text)*7), 43-(d.CurrentConversation.HighlightedMessage.Index-d.CurrentConversation.Messages[i].Index)*12, d.CurrentConversation.Messages[i].Text)
				}
			} else if d.CurrentConversation.Messages[i].Index > d.CurrentConversation.HighlightedMessage.Index {
				if d.CurrentConversation.Messages[i].Person != *d.SelfIdentity {
					drawText(img, 0, 43+(d.CurrentConversation.Messages[i].Index-d.CurrentConversation.HighlightedMessage.Index)*12, d.CurrentConversation.Messages[i].Text)
				} else {
					drawText(img, dimensions.Dx()-(len(d.CurrentConversation.Messages[i].Text)*7), 43+(d.CurrentConversation.Messages[i].Index-d.CurrentConversation.HighlightedMessage.Index)*12, d.CurrentConversation.Messages[i].Text)
				}
			}
		}
	}

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
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

// drawHLine draws a white horizontal line from one X location to another. x2 has to be greater than x1.
func drawHLine(img *image.RGBA, x1 int, y int, x2 int) {
	col := color.RGBA{255, 255, 255, 255}
	drawHLineCol(img, x1, y, x2, col)
}

// drawHLineCol draws a horizontal line in a color of your choice from one X location to another. x2 has to be greater than x1.
func drawHLineCol(img *image.RGBA, x1 int, y int, x2 int, col color.RGBA) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// drawVLine draws a verticle line from one Y location to another. y2 has to be greater than y1.
func drawVLine(img *image.RGBA, y1 int, x int, y2 int) {
	col := color.RGBA{255, 255, 255, 255}
	drawVLineCol(img, y1, x, y2, col)
}

// drawVLineCol draws a vertical line in a color of your choice from one Y location to another. y2 has to be greater than y1.
func drawVLineCol(img *image.RGBA, y1 int, x int, y2 int, col color.RGBA) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// drawBlackFilledBox draws a filled blacck box from one X and Y location to another.
func drawBlackFilledBox(img *image.RGBA, x1 int, y1 int, x2 int, y2 int) {
	col := color.RGBA{0, 0, 0, 255}
	for ; y1 <= y2; y1++ {
		drawHLineCol(img, x1, y1, x2, col)
	}
}
