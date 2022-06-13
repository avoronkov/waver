package signals

type Input interface {
	Run(chan<- Interface) error
	Close() error
}
