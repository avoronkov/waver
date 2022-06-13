package signals

type Tempo struct {
	Tempo int
}

func (*Tempo) SignalType() string {
	return "tempo"
}
