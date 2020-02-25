package main

// Modes holds single instances of each mode.
type Modes struct {
	// Mode is the active mode.
	Mode
	Command Mode
	Input   Mode
	Select  Mode
	Switch  Mode
}

// SwitchMode switches active mode.
func (modes *Modes) SwitchMode(mode Mode) {
	if mode == modes.Mode {
		return
	}
	if modes.Mode != nil {
		modes.Mode.Hide()
	}
	modes.Mode = mode
	modes.Mode.Show()
}
