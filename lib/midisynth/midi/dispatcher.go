package midi

type Dispatcher interface {
	Up(channel int, knob int)
	Down(channel int, knob int)
}
