package notes

import "fmt"

// 19th root of 2.0
const Edo19Step = 1.0371550444462

type Edo19 struct {
	notes   map[int]map[string]Note
	reprs   map[string]Note
	indexes map[int]Note
}

var _ (Scale) = (*Edo19)(nil)
var _ (EdoScale) = (*Edo19)(nil)

func NewEdo19() *Edo19 {
	e := &Edo19{}
	e.init()
	return e
}

func (e *Edo19) init() {
	e.notes = make(map[int]map[string]Note)
	e.reprs = make(map[string]Note)
	e.indexes = make(map[int]Note)
	afreq := 55.0
	num := 1
	for oct := 1; oct <= 7; oct++ {
		num = e.buildOctaveRepr(oct, afreq, num)
		afreq *= 2.0
	}
}

func (s *Edo19) buildOctaveRepr(oct int, afreq float64, num int) int {
	s.notes[oct] = make(map[string]Note)

	var idx int

	a := afreq
	idx = num + 14
	nt := Note{idx, a}
	s.reprs[fmt.Sprintf("A%d", oct)] = nt
	s.notes[oct]["A"] = nt
	s.indexes[idx] = nt

	as := a * Edo19Step
	idx = num + 15
	nt = Note{idx, as}
	s.reprs[fmt.Sprintf("As%d", oct)] = nt
	s.notes[oct]["a"] = nt
	s.indexes[idx] = nt

	bb := as * Edo19Step
	idx = num + 16
	nt = Note{idx, bb}
	s.reprs[fmt.Sprintf("Bb%d", oct)] = nt
	s.notes[oct]["9"] = nt
	s.indexes[idx] = nt

	b := bb * Edo19Step
	idx = num + 17
	nt = Note{idx, b}
	s.reprs[fmt.Sprintf("B%d", oct)] = nt
	s.notes[oct]["B"] = nt
	s.indexes[idx] = nt

	bs := b * Edo19Step
	idx = num + 18
	nt = Note{idx, bs}
	s.reprs[fmt.Sprintf("Bs%d", oct)] = nt
	s.notes[oct]["b"] = nt
	s.indexes[idx] = nt

	ab := a / Edo19Step
	idx = num + 13
	nt = Note{idx, ab}
	s.reprs[fmt.Sprintf("Ab%d", oct)] = nt
	s.notes[oct]["8"] = nt
	s.indexes[idx] = nt

	gs := ab / Edo19Step
	idx = num + 12
	nt = Note{idx, gs}
	s.reprs[fmt.Sprintf("Gs%d", oct)] = nt
	s.notes[oct]["g"] = nt
	s.indexes[idx] = nt

	g := gs / Edo19Step
	idx = num + 11
	nt = Note{idx, g}
	s.reprs[fmt.Sprintf("G%d", oct)] = nt
	s.notes[oct]["G"] = nt
	s.indexes[idx] = nt

	gb := g / Edo19Step
	idx = num + 10
	nt = Note{idx, gb}
	s.reprs[fmt.Sprintf("Gb%d", oct)] = nt
	s.notes[oct]["7"] = nt
	s.indexes[idx] = nt

	fs := gb / Edo19Step
	idx = num + 9
	nt = Note{idx, fs}
	s.reprs[fmt.Sprintf("Fs%d", oct)] = nt
	s.notes[oct]["f"] = nt
	s.indexes[idx] = nt

	f := fs / Edo19Step
	idx = num + 8
	nt = Note{idx, f}
	s.reprs[fmt.Sprintf("F%d", oct)] = nt
	s.notes[oct]["F"] = nt
	s.indexes[idx] = nt

	fb := f / Edo19Step
	idx = num + 7
	nt = Note{idx, fb}
	s.reprs[fmt.Sprintf("Fb%d", oct)] = nt
	s.reprs[fmt.Sprintf("Es%d", oct)] = nt
	s.notes[oct]["e"] = nt
	s.indexes[idx] = nt

	e := fb / Edo19Step
	idx = num + 6
	nt = Note{idx, e}
	s.reprs[fmt.Sprintf("E%d", oct)] = nt
	s.notes[oct]["E"] = nt
	s.indexes[idx] = nt

	eb := e / Edo19Step
	idx = num + 5
	nt = Note{idx, eb}
	s.reprs[fmt.Sprintf("Eb%d", oct)] = nt
	s.notes[oct]["6"] = nt
	s.indexes[idx] = nt

	ds := eb / Edo19Step
	idx = num + 4
	nt = Note{idx, ds}
	s.reprs[fmt.Sprintf("Ds%d", oct)] = nt
	s.notes[oct]["d"] = nt
	s.indexes[idx] = nt

	d := ds / Edo19Step
	idx = num + 3
	nt = Note{idx, d}
	s.reprs[fmt.Sprintf("D%d", oct)] = nt
	s.notes[oct]["D"] = nt
	s.indexes[idx] = nt

	db := d / Edo19Step
	idx = num + 2
	nt = Note{idx, db}
	s.reprs[fmt.Sprintf("Db%d", oct)] = nt
	s.notes[oct]["5"] = nt
	s.indexes[idx] = nt

	cs := db / Edo19Step
	idx = num + 1
	nt = Note{idx, cs}
	s.reprs[fmt.Sprintf("Cs%d", oct)] = nt
	s.notes[oct]["c"] = nt
	s.indexes[idx] = nt

	c := cs / Edo19Step
	idx = num
	nt = Note{idx, c}
	s.reprs[fmt.Sprintf("C%d", oct)] = nt
	s.notes[oct]["C"] = nt
	s.indexes[idx] = nt

	/*
		cb := c / Edo19Step
		nt = Note{num, cb}
		s.reprs[fmt.Sprintf("Cb%d", oct)] = nt
		s.indexes[num] = nt
	*/

	return num + 19
}

func (e *Edo19) Note(octave int, note string) (Note, bool) {
	hz, ok := e.notes[octave][note]
	return Note(hz), ok
}

func (e *Edo19) Parse(str string) (Note, bool) {
	n, ok := e.reprs[str]
	return n, ok
}

func (e *Edo19) ByNumber(n int) (Note, bool) {
	nt, ok := e.indexes[n]
	return nt, ok
}

func (e *Edo19) Edo() int {
	return 19
}
