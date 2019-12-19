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
func (ms *Modes) SwitchMode(m Mode) {
	if ms.Mode == m {
		return
	}
	if ms.Mode != nil {
		ms.Mode.Hide()
	}
	ms.Mode = m
	ms.Mode.Show()
}
