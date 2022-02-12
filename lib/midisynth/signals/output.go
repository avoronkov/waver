package signals

type Output interface {
	ProcessAsync(*Signal)
	Close() error
}
