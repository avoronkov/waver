package signals

type Input interface {
	Run(chan<- *Signal) error
	Close() error
}
