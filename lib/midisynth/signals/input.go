package signals

type Input interface {
	Start(chan<- *Signal) error
	Close() error
}
