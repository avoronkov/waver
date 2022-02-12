package signals

type Output interface {
	ProcessAsync(time float64, sig *Signal)
	Close() error
}
