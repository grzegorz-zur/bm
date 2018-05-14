package bm

type Modes struct {
	Mode
	Command Mode
	Input   Mode
	Switch  Mode
}

func (modes *Modes) SwitchMode(mode Mode) {
	if modes.Mode == mode {
		return
	}
	if modes.Mode != nil {
		modes.Mode.Hide()
	}
	modes.Mode = mode
	modes.Mode.Show()
}
