package picodoomsdaymessenger

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Device is the main structure that holds all the information about the device. It has a State, a StateHistory, and an LEDAnimation.
type Device struct {
	State                    *State
	StateHistory             []*State
	LEDAnimation             *LEDAnimation
	Conversations            []*Conversation
	CurrentConversationIndex int
	SelfIdentity             Person
	CurrentKeyboardButton    *KeyboardButton
	SendUsingRadio           func(packet []byte) (err error)
}

type KeyboardButton struct {
	Characters            []string
	LastPress             time.Time
	CurrentCharacterIndex int
}

// Conversation is a conversation with a person. It contains a list of Messages and a Person that the conversation is with.
type Conversation struct {
	Messages                []Message
	HighlightedMessageIndex int
	KeyboardBuffer          string
	Name                    string
	People                  []Person
}

// Person is a representation of another device. A Person has a name and a unique identifier
type Person struct {
	Name string
	ID   int
}

// Message is a message sent inside a Conversation. It contains the time it was sent, the time it was recieved and the content of the message.
type Message struct {
	Text     string
	Person   Person
	TimeSent time.Time
}

// State is the current state of the device. It contains all the information about what is currently being displayed.
type State struct {
	Title                string
	Content              []MenuItem
	HighlightedItemIndex int
	LoadAction           func(d *Device) (err error)
}

// MenuItem is a structure that holds data that can be displayed on the screen. It contains a title and an action that is run when the item is selected.
type MenuItem struct {
	Text          string
	Action        func(d *Device) (err error)
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

// Define errors
var (
	ErrCursorIconBoxBoolTypeError         = errors.New("data is not a bool")
	ErrDefaultMenuItem                    = errors.New("this error should not appear, if it does please tell us")
	ErrRadioSendNotDefined                = errors.New("radio send function not defined by user")
	ErrConversationReaderAcceptDisallowed = errors.New("cannot accept in conversation reader")
	ErrGoBackStateRootState               = errors.New("already at root state")
	ErrInvalidMessage                     = errors.New("invalid message, prefix incorrect")
)

// Define the Keyboard Buttons
var (
	KeyboardButton1 = &KeyboardButton{[]string{"1", "2"}, time.Time{}, 0}
	KeyboardButton2 = &KeyboardButton{[]string{"a", "b", "c"}, time.Time{}, 0}
	KeyboardButton3 = &KeyboardButton{[]string{"d", "e", "f"}, time.Time{}, 0}
	KeyboardButton4 = &KeyboardButton{[]string{"g", "h", "i"}, time.Time{}, 0}
	KeyboardButton5 = &KeyboardButton{[]string{"j", "k", "l"}, time.Time{}, 0}
	KeyboardButton6 = &KeyboardButton{[]string{"m", "n", "o"}, time.Time{}, 0}
	KeyboardButton7 = &KeyboardButton{[]string{"p", "q", "r", "s"}, time.Time{}, 0}
	KeyboardButton8 = &KeyboardButton{[]string{"t", "u", "v"}, time.Time{}, 0}
	KeyboardButton9 = &KeyboardButton{[]string{"w", "x", "y", "z"}, time.Time{}, 0}
	KeyboardButton0 = &KeyboardButton{[]string{" "}, time.Time{}, 0}
)

// Define default People

// PersonYou is a default person that is used for your self identity. Do not use this to identify yourself, use d.SelfIdentity instead.
var PersonYou = Person{"You", 0}

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
			return ErrCursorIconBoxBoolTypeError
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
			return ErrDefaultMenuItem
		},
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
		CursorIcon: CursorIconLeftArrow,
	}

	// Main Menu Items

	// MainMenuItemConversations is a MenuItem that goes to the StateConversations menu.
	MainMenuItemConversations MenuItem = MenuItem{
		Text: "Conversations",
		Action: func(d *Device) (err error) {
			err = d.ChangeStateWithHistory(&StateConversationsMenu)
			if err != nil {
				return err
			}
			return nil
		},

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

		GetCursorData: func(d *Device) (data any, err error) {
			return d.LEDAnimation == &LEDAnimationSOS, nil
		},
		CursorIcon: CursorIconBox,
	}
	// Conversation Menu Items
	ConversationsMenuItemNew MenuItem = MenuItem{
		Text: "New Conversation",
		Action: func(d *Device) (err error) {
			d.ChangeStateWithHistory(&StateNewConversation)
			return nil
		},
		CursorIcon: CursorIconRightArrow,
	}
)

// Define States
var (
	// StateDefault is a State that does nothing. It is used as a placeholder for the default State of the Device.
	StateDefault = State{
		Title:                "DefaultState",
		Content:              []MenuItem{MenuItemDefault},
		HighlightedItemIndex: 0,
	}
	// StateConversationReader is a special State that is used when reading a Conversation.
	StateConversationReader = State{
		Title:   "",
		Content: []MenuItem{},
	}
	// StateMainMenu is a State that shows the main menu.
	StateMainMenu = State{
		Title:                "Main Menu",
		Content:              []MenuItem{MainMenuItemConversations, MainMenuItemPeople, MainMenuItemGames, MainMenuItemDemos, MainMenuItemTools, MainMenuItemSettings},
		HighlightedItemIndex: 0,
	}
	// StateConversationsMenu is a State that shows the conversations menu.
	StateConversationsMenu = State{
		Title:                "Conversations",
		Content:              []MenuItem{GlobalMenuItemGoBack, ConversationsMenuItemNew},
		HighlightedItemIndex: 0,
	}
	// StateConversationsMenuOld is a copy of StateConversationsMenu that can be used as a starting point to reset StateConversationsMenu.
	StateConversationsMenuOld = StateConversationsMenu
	// StateNewConversation is a special State that is used when creating a new Conversation.
	StateNewConversation = State{
		Title:                "New Conversation",
		Content:              []MenuItem{GlobalMenuItemGoBack},
		HighlightedItemIndex: 0,
	}
	// StatePeopleMenu is a State that shows the people menu.
	StatePeopleMenu = State{
		Title:                "People",
		Content:              []MenuItem{GlobalMenuItemGoBack},
		HighlightedItemIndex: 0,
	}
	// StateGamesMenu is a State that shows the games menu.
	StateGamesMenu = State{
		Title:                "Games",
		Content:              []MenuItem{GlobalMenuItemGoBack},
		HighlightedItemIndex: 0,
	}
	// StateDemosMenu is a State that shows the demos menu.
	StateDemosMenu = State{
		Title:                "Demos",
		Content:              []MenuItem{GlobalMenuItemGoBack, DemoMenuItemRGB},
		HighlightedItemIndex: 0,
	}
	// StateToolsMenu is a State that shows the tools menu.
	StateToolsMenu = State{
		Title:                "Tools",
		Content:              []MenuItem{GlobalMenuItemGoBack, ToolsMenuItemSOS},
		HighlightedItemIndex: 0,
	}
	// StateSettingsMenu is a State that shows the settings menu.
	StateSettingsMenu = State{
		Title:                "Settings",
		Content:              []MenuItem{GlobalMenuItemGoBack},
		HighlightedItemIndex: 0,
	}
)

// Define LED animations. They are made of multiple frames of 6 colors.
var (
	// LEDAnimationDefault is the default LED animation. It is used when no other animation is active and is simply black.
	LEDAnimationDefault = LEDAnimation{
		FrameDuration: 100 * time.Millisecond,
		CurrentFrame:  0,
		Frames: [][6]color.RGBA{
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
func NewDevice() (d *Device, err error) {
	rand.Seed(time.Now().UnixNano())
	PersonYou.ID = rand.Intn(2147483647) // Max value of an int32
	return &Device{
		State:                    &StateMainMenu,
		StateHistory:             []*State{&StateMainMenu},
		LEDAnimation:             &LEDAnimationDefault,
		Conversations:            []*Conversation{},
		SelfIdentity:             PersonYou,
		CurrentConversationIndex: 0,
		CurrentKeyboardButton:    KeyboardButton0,
		SendUsingRadio: func(packet []byte) (err error) {
			return ErrRadioSendNotDefined
		},
	}, nil
}

// RecieveFromRadio takes in the payload of a radio packet, usually recieved from the RFM9x radio.
func (d *Device) ReceiveFromRadio(packetPayload []byte) (err error) {
	payloadMessage, err := d.BytesToMessage(packetPayload)
	if err != nil {
		return err
	}

	newConversation := d.NewConversation(payloadMessage.Person)
	newConversation.Messages = append(newConversation.Messages, payloadMessage)
	newConversation.HighlightedMessageIndex = len(newConversation.Messages) - 1
	newConversation.Name = fmt.Sprint(payloadMessage.Person.ID)

	d.UpdateConversationsMenu()
	return nil
}

// NewConversation creates a blank new Conversation with a person and adds it to the Device. It also returns a pointer to that Conversation.
func (d *Device) NewConversation(p Person) (c *Conversation) {
	newConversation := &Conversation{People: []Person{d.SelfIdentity, p}}
	d.Conversations = append(d.Conversations, newConversation)
	return newConversation
}

func (d *Device) UpdateConversationsMenu() {
	StateConversationsMenu = StateConversationsMenuOld
	for i := 0; i < len(d.Conversations); i++ {
		// Define a seperate variable to seperate the increasing i from the functions defined here.
		j := i
		StateConversationsMenu.Content = append(StateConversationsMenu.Content, MenuItem{
			Text: d.Conversations[j].Name,
			Action: func(d *Device) (err error) {
				d.CurrentConversationIndex = j
				err = d.ChangeStateWithHistory(&StateConversationReader)
				return err
			},
			CursorIcon: CursorIconRightArrow,
		})
	}
	StateConversationsMenu.HighlightedItemIndex = len(StateConversationsMenu.Content) - 1
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
		return ErrGoBackStateRootState
	}
	err = d.ChangeStateWithoutHistory(d.StateHistory[len(d.StateHistory)-2])
	d.StateHistory = d.StateHistory[0 : len(d.StateHistory)-1]
	return err
}

// InputEvent is a string that represents a button press.
type InputEvent string

const (
	InputEventUp                InputEvent = "up"
	InputEventDown              InputEvent = "down"
	InputEventLeft              InputEvent = "left"
	InputEventRight             InputEvent = "right"
	InputEventAccept            InputEvent = "accept"
	InputEventFunction1         InputEvent = "function1"
	InputEventFunction2         InputEvent = "function2"
	InputEventFunction3         InputEvent = "function3"
	InputEventFunction4         InputEvent = "function4"
	InputEventOpenSettings      InputEvent = "openSettings"
	InputEventOpenPeople        InputEvent = "openPeople"
	InputEventOpenConversations InputEvent = "openConversations"
	InputEventOpenMainMenu      InputEvent = "openMainMenu"
	InputEventNumber1           InputEvent = "number1"
	InputEventNumber2           InputEvent = "number2"
	InputEventNumber3           InputEvent = "number3"
	InputEventNumber4           InputEvent = "number4"
	InputEventNumber5           InputEvent = "number5"
	InputEventNumber6           InputEvent = "number6"
	InputEventNumber7           InputEvent = "number7"
	InputEventNumber8           InputEvent = "number8"
	InputEventNumber9           InputEvent = "number9"
	InputEventNumber0           InputEvent = "number0"
	InputEventStar              InputEvent = "star"
	InputEventPound             InputEvent = "pound"
)

// ProcessInputEvent will take in an InputEvent and run appropriate actions based on the event.
func (d *Device) ProcessInputEvent(inputEvent InputEvent) (err error) {
	// Process the keys that are always available.
	switch inputEvent {
	case InputEventUp:
		{
			err = d.ProcessInputEventUp()
			return err
		}
	case InputEventDown:
		{
			err = d.ProcessInputEventDown()
			return err
		}
	case InputEventAccept:
		{
			err = d.ProcessInputEventAccept()
			return err
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
	case InputEventOpenConversations:
		{
			err = d.ChangeStateWithHistory(&StateConversationsMenu)
			return err
		}
	case InputEventOpenMainMenu:
		{
			err = d.ChangeStateWithHistory(&StateMainMenu)
			return err
		}
	}
	// Process the keys that are available in the conversationreader state.
	if d.State == &StateConversationReader {
		switch inputEvent {
		case InputEventNumber1:
			{
				err = d.ProcessConversationInputEventNumber1()
				return err
			}
		case InputEventNumber2:
			{
				err = d.ProcessConversationInputEventNumber2()
				return err
			}
		case InputEventNumber3:
			{
				err = d.ProcessConversationInputEventNumber3()
				return err
			}
		case InputEventNumber4:
			{
				err = d.ProcessConversationInputEventNumber4()
				return err
			}
		case InputEventNumber5:
			{
				err = d.ProcessConversationInputEventNumber5()
				return err
			}
		case InputEventNumber6:
			{
				err = d.ProcessConversationInputEventNumber6()
				return err
			}
		case InputEventNumber7:
			{
				err = d.ProcessConversationInputEventNumber7()
				return err
			}
		case InputEventNumber8:
			{
				err = d.ProcessConversationInputEventNumber8()
				return err
			}
		case InputEventNumber9:
			{
				err = d.ProcessConversationInputEventNumber9()
				return err
			}
		case InputEventNumber0:
			{
				err = d.ProcessConversationInputEventNumber0()
				return err
			}
		}
	}
	return nil
}

func (d *Device) ProcessInputEventUp() (err error) {
	if d.State != &StateConversationReader {
		if d.State.HighlightedItemIndex <= 0 {
			d.State.HighlightedItemIndex = len(d.State.Content) - 1
		} else {
			d.State.HighlightedItemIndex--
		}
	} else {
		if d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex <= 0 {
			d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex = len(d.Conversations[d.CurrentConversationIndex].Messages) - 1
		} else {
			d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex = d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex - 1
		}
	}
	return nil
}

func (d *Device) ProcessInputEventDown() (err error) {
	if d.State != &StateConversationReader {
		if d.State.HighlightedItemIndex >= len(d.State.Content)-1 {
			d.State.HighlightedItemIndex = 0
		} else {
			d.State.HighlightedItemIndex++
		}
	} else {
		if d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex >= len(d.Conversations[d.CurrentConversationIndex].Messages)-1 {
			d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex = 0
		} else {
			d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex++
		}
	}
	return nil
}

func (d *Device) ProcessInputEventAccept() (err error) {
	if d.State != &StateConversationReader {
		err = d.State.Content[d.State.HighlightedItemIndex].Action(d)
		return err
	}
	packetToSend, err := d.MesageToBytes(Message{
		Text:   d.Conversations[d.CurrentConversationIndex].KeyboardBuffer + d.CurrentKeyboardButton.Characters[d.CurrentKeyboardButton.CurrentCharacterIndex],
		Person: d.SelfIdentity,
	})
	if err != nil {
		return err
	}
	d.Conversations[d.CurrentConversationIndex].KeyboardBuffer = ""
	d.CurrentKeyboardButton = &KeyboardButton{Characters: []string{""}, CurrentCharacterIndex: 0}
	return d.SendUsingRadio(packetToSend)
}

func (d *Device) ProcessConversationInputEventNumber1() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton1)
}

func (d *Device) ProcessConversationInputEventNumber2() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton2)
}

func (d *Device) ProcessConversationInputEventNumber3() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton3)
}

func (d *Device) ProcessConversationInputEventNumber4() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton4)
}

func (d *Device) ProcessConversationInputEventNumber5() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton5)
}

func (d *Device) ProcessConversationInputEventNumber6() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton6)
}

func (d *Device) ProcessConversationInputEventNumber7() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton7)
}

func (d *Device) ProcessConversationInputEventNumber8() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton8)
}

func (d *Device) ProcessConversationInputEventNumber9() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton9)
}

func (d *Device) ProcessConversationInputEventNumber0() (err error) {
	return d.ProcessConversationInputEventNumber(KeyboardButton0)
}

func (d *Device) ProcessConversationInputEventNumber(button *KeyboardButton) (err error) {
	if d.CurrentKeyboardButton != button {
		d.Conversations[d.CurrentConversationIndex].KeyboardBuffer += d.CurrentKeyboardButton.Characters[d.CurrentKeyboardButton.CurrentCharacterIndex]
		d.CurrentKeyboardButton = button
		d.CurrentKeyboardButton.CurrentCharacterIndex = 0
	} else {
		if d.CurrentKeyboardButton.CurrentCharacterIndex >= len(d.CurrentKeyboardButton.Characters)-1 {
			d.CurrentKeyboardButton.CurrentCharacterIndex = 0
		} else {
			d.CurrentKeyboardButton.CurrentCharacterIndex++
		}
	}
	return nil
}

// MesageToBytes converts a Message to a compressed byte array.
func (d *Device) MesageToBytes(input Message) (output []byte, err error) {
	staringBytes := []byte{0x64, 0x6F, 0x6F, 0x6D} // ASCII for "doom"
	seperatorByte := byte(0xcc)
	bytesToSend := make([]byte, 0)
	bytesToSend = append(bytesToSend, staringBytes...)
	bytesToSend = append(bytesToSend, []byte(fmt.Sprint(input.Person.ID))...)
	bytesToSend = append(bytesToSend, seperatorByte)
	bytesToSend = append(bytesToSend, []byte(input.Person.Name)...)
	bytesToSend = append(bytesToSend, seperatorByte)
	bytesToSend = append(bytesToSend, []byte(input.Text)...)
	return bytesToSend, nil
}

// BytesToMessage converts a compressed byte array to a Message.
func (d *Device) BytesToMessage(input []byte) (output Message, err error) {
	startingBytes := []byte{0x64, 0x6F, 0x6F, 0x6D} // ASCII for "doom"
	if !bytes.HasPrefix(input, startingBytes) {
		return output, ErrInvalidMessage
	}
	seperatorByte := byte(0xcc)
	receivedBytesSplit := bytes.Split(input, []byte{seperatorByte})
	personID := receivedBytesSplit[0]
	output.Person.ID = int(personID[0])
	output.Person.Name = string(receivedBytesSplit[1])
	output.Text = string(receivedBytesSplit[2])
	return output, nil
}

// GetFrame will take in a Device and return an image based on the state.
func GetFrame(dimensions image.Rectangle, d *Device) (frame image.Image, err error) {
	img := image.NewRGBA(dimensions)

	if d.State != &StateConversationReader && d.State != &StateNewConversation {
		// Draw the content with the currently highlighted item in the middle of the screen and the other items above and below it.
		for i := 0; i < len(d.State.Content); i++ {
			if i == d.State.HighlightedItemIndex {
				drawText(img, 0, 43, d.State.Content[i].Text)
			} else if i < d.State.HighlightedItemIndex {
				drawText(img, 0, 43-(d.State.HighlightedItemIndex-i)*12, d.State.Content[i].Text)
			} else if i > d.State.HighlightedItemIndex {
				drawText(img, 0, 43+(i-d.State.HighlightedItemIndex)*12, d.State.Content[i].Text)
			}
		}

		// Draw the title.
		drawBlackFilledBox(img, 0, 0, dimensions.Dx(), 16)
		drawText(img, 0, 13, d.State.Title)
		drawHLine(img, 0, 15, dimensions.Dx())

		// Draw the cursor. If the cursor is a checkbox, check if the checkbox is checked or not.
		var cursorData any
		if d.State.Content[d.State.HighlightedItemIndex].GetCursorData != nil {
			cursorData, err = d.State.Content[d.State.HighlightedItemIndex].GetCursorData(d)
			if err != nil {
				return nil, err
			}
		}
		err = d.State.Content[d.State.HighlightedItemIndex].CursorIcon(img, dimensions.Dx()-7, 36, cursorData)
		if err != nil {
			return nil, err
		}
	} else if d.State == &StateConversationReader {
		// Draw the conversation with the most recent message at the bottom of the screen.
		for i := 0; i < len(d.Conversations[d.CurrentConversationIndex].Messages); i++ {
			if i == d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex {
				if d.Conversations[d.CurrentConversationIndex].Messages[i].Person != d.SelfIdentity {
					drawText(img, 0, 43, "> "+d.Conversations[d.CurrentConversationIndex].Messages[i].Text)
				} else {
					drawText(img, dimensions.Dx()-((len(d.Conversations[d.CurrentConversationIndex].Messages[i].Text)+2)*7), 43, d.Conversations[d.CurrentConversationIndex].Messages[i].Text+" <")
				}
			} else if i < d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex {
				if d.Conversations[d.CurrentConversationIndex].Messages[i].Person != d.SelfIdentity {
					drawText(img, 0, 43-(d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex-i)*12, "> "+d.Conversations[d.CurrentConversationIndex].Messages[i].Text)
				} else {
					drawText(img, dimensions.Dx()-((len(d.Conversations[d.CurrentConversationIndex].Messages[i].Text)+2)*7), 43-(d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex-i)*12, d.Conversations[d.CurrentConversationIndex].Messages[i].Text+" <")
				}
			} else if i > d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex {
				if d.Conversations[d.CurrentConversationIndex].Messages[i].Person != d.SelfIdentity {
					drawText(img, 0, 43+(i-d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex)*12, "> "+d.Conversations[d.CurrentConversationIndex].Messages[i].Text)
				} else {
					drawText(img, dimensions.Dx()-((len(d.Conversations[d.CurrentConversationIndex].Messages[i].Text)+2)*7), 43+(i-d.Conversations[d.CurrentConversationIndex].HighlightedMessageIndex)*12, d.Conversations[d.CurrentConversationIndex].Messages[i].Text+" <")
				}
			}
		}
		drawBlackFilledBox(img, 0, 0, dimensions.Dx(), 16)
		drawText(img, 0, 13, d.Conversations[d.CurrentConversationIndex].Name)
		drawHLine(img, 0, 15, dimensions.Dx())
		drawBlackFilledBox(img, 0, ((dimensions.Dy()*75)/100)-1, dimensions.Dx(), dimensions.Dy())
		drawHLine(img, 0, (dimensions.Dy()*75)/100, dimensions.Dx())
		drawText(img, 0, (dimensions.Dy()*75)/100+13, d.Conversations[d.CurrentConversationIndex].KeyboardBuffer+d.CurrentKeyboardButton.Characters[d.CurrentKeyboardButton.CurrentCharacterIndex])
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
