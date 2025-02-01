package notes

type Scale interface {
	Note(octave int, note string) (Note, bool)
	Parse(s string) (Note, bool)
	ByNumber(n int) (Note, bool)
}

type EdoScale interface {
	Edo() int
}

type StdFuncsScale interface {
	Std() []byte
}
