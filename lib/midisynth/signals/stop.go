package signals

type Stop struct{}

func (*Stop) SignalType() string {
	return "stop"
}
