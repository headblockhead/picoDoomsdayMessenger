package picodoomsdaymessenger

import (
	"errors"
	"testing"
)

func TestDefaults(t *testing.T) {
	// Create a new Machine
	device := NewDevice()

	// Test the default state
	if device.State != &StateMainMenu {
		t.Errorf("The default state should be StateMainMenu but is %v", device.State)
	}

	// Test the default state history
	if len(device.StateHistory) != 1 {
		t.Errorf("The default state history should only contain 1 item but it contains %v", device.StateHistory)
	}

	// Test the default LED animation
	if device.LEDAnimation != &LEDAnimationDefault {
		t.Errorf("The default LED animation should be LEDAnimationDefault but is %v", device.LEDAnimation)
	}
}

func TestChangeLEDAnimationWithoutContinue(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testLEDAnimation1 := LEDAnimation{
		CurrentFrame: 0,
	}
	testLEDAnimation2 := LEDAnimation{
		CurrentFrame: 1,
	}
	device.LEDAnimation = &testLEDAnimation1

	device.ChangeLEDAnimationWithoutContinue(&testLEDAnimation2)

	if device.LEDAnimation != &testLEDAnimation2 {
		t.Errorf("The LEDAnimation should be testLEDAnimation2 but is %v", device.LEDAnimation)
	}
	if device.LEDAnimation.CurrentFrame != 0 {
		t.Errorf("The LEDAnimation.CurrentFrame should be 0 but is %v", device.LEDAnimation.CurrentFrame)
	}
}

func TestChangeLEDAnimationWithContinue(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testLEDAnimation1 := LEDAnimation{
		CurrentFrame: 0,
	}
	testLEDAnimation2 := LEDAnimation{
		CurrentFrame: 1,
	}
	device.LEDAnimation = &testLEDAnimation1

	device.ChangeLEDAnimationWithContinue(&testLEDAnimation2)

	if device.LEDAnimation != &testLEDAnimation2 {
		t.Errorf("The LEDAnimation should be testLEDAnimation2 but is %v", device.LEDAnimation)
	}
	if device.LEDAnimation.CurrentFrame != 1 {
		t.Errorf("The LEDAnimation.CurrentFrame should be 1 but is %v", device.LEDAnimation.CurrentFrame)
	}
}

func TestChangeStateWithHistory(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	testState1 := State{
		HighlightedItem: &MenuItemDefault,
		LoadAction: func(d *Device) (err error) {
			return errors.New("test error")
		},
	}
	device.State = &testState0

	err := device.ChangeStateWithHistory(&testState1)

	if err.Error() != "test error" {
		t.Errorf("The error should be \"test error\" but is %v", err)
	}
	if device.State != &testState1 {
		t.Errorf("The state should be testState1 but is %v", device.State)
	}
	if len(device.StateHistory) != 2 {
		t.Errorf("The state history should contain 2 items but contains %v", device.StateHistory)
	}
}

func TestChangeStateWithoutHistory(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	testState1 := State{
		HighlightedItem: &MenuItemDefault,
	}
	device.State = &testState0

	err := device.ChangeStateWithoutHistory(&testState1)

	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State != &testState1 {
		t.Errorf("The state should be testState1 but is %v", device.State)
	}
	if len(device.StateHistory) != 1 {
		t.Errorf("The state history should contain 1 item but contains %v", device.StateHistory)
	}
}

func TestGoBackState(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	testState1 := State{
		HighlightedItem: &MenuItemDefault,
	}
	device.StateHistory = []*State{&testState0, &testState1}
	device.State = &testState1

	err := device.GoBackState()
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}

	if device.State != &testState0 {
		t.Errorf("The state should be testState0 but is %v", device.State)
	}

	device.StateHistory = []*State{&testState0}
	device.State = &testState0

	err = device.GoBackState()
	if err == nil {
		t.Errorf("The error should not be nil")
	}
}

func TestProcessInputEventUp(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	menuItem0 := MenuItem{Text: "test0", Index: 0}
	menuItem1 := MenuItem{Text: "test1", Index: 1}
	testState0 := State{
		Content:         []MenuItem{menuItem0, menuItem1},
		HighlightedItem: &menuItem1,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventUp)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State.HighlightedItem.Text != "test0" {
		t.Errorf("The highlighted item should be menuItem0 but is %v", device.State.HighlightedItem)
	}

	device.State.HighlightedItem = &menuItem0
	err = device.ProcessInputEvent(InputEventUp)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State.HighlightedItem.Text != "test1" {
		t.Errorf("The highlighted item should be menuItem1 but is %v", device.State.HighlightedItem)
	}
}

func TestProcessInputEventDown(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	menuItem0 := MenuItem{Text: "test0", Index: 0}
	menuItem1 := MenuItem{Text: "test1", Index: 1}
	testState0 := State{
		Content:         []MenuItem{menuItem0, menuItem1},
		HighlightedItem: &menuItem0,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventDown)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State.HighlightedItem.Text != "test1" {
		t.Errorf("The highlighted item should be menuItem1 but is %v", device.State.HighlightedItem)
	}

	device.State.HighlightedItem = &menuItem1
	err = device.ProcessInputEvent(InputEventDown)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State.HighlightedItem.Text != "test0" {
		t.Errorf("The highlighted item should be menuItem0 but is %v", device.State.HighlightedItem)
	}
}

func TestProcessInputEventAccept(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	menuItem0 := MenuItem{Text: "test0", Index: 0, Action: func(d *Device) (err error) {
		return errors.New("test error")
	}}
	testState0 := State{
		Content:         []MenuItem{menuItem0},
		HighlightedItem: &menuItem0,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventAccept)
	if err.Error() != "test error" {
		t.Errorf("The error should be \"test error\" but is %v", err)
	}
}

func TestProcessInputEventOpenSettings(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventOpenSettings)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State != &StateSettingsMenu {
		t.Errorf("The state should be StateSettings but is %v", device.State)
	}
}

func TestProcessInputEventOpenPeople(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventOpenPeople)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State != &StatePeopleMenu {
		t.Errorf("The state should be StatePeople but is %v", device.State)
	}
}

func TestProcessInputEventOpenMessages(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventOpenMessages)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State != &StateMessagesMenu {
		t.Errorf("The state should be StateMessages but is %v", device.State)
	}
}

func TestProcessInputEventOpenMainMenu(t *testing.T) {
	// Create a new Machine
	device := NewDevice()
	testState0 := State{
		HighlightedItem: &MenuItemDefault,
	}
	device.State = &testState0

	err := device.ProcessInputEvent(InputEventOpenMainMenu)
	if err != nil {
		t.Errorf("The error should be nil but is %v", err)
	}
	if device.State != &StateMainMenu {
		t.Errorf("The state should be StateMainMenu but is %v", device.State)
	}
}