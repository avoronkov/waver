package notes

type Scale interface {
	Note(octave int, note string) (hz float64, ok bool)
}

type EdoScale interface {
	Edo() int
}
