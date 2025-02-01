package notes

import (
	"fmt"

	"github.com/avoronkov/waver/etc/std"
)

// 12th root of 2.0
const HalfStep = 1.0594630943593

// Standard
type Standard struct {
	notes   map[int]map[string]Note
	reprs   map[string]Note
	indexes map[int]Note
}

var _ Scale = (*Standard)(nil)
var _ EdoScale = (*Standard)(nil)
var _ StdFuncsScale = (*Standard)(nil)

func NewStandard() *Standard {
	s := &Standard{}
	s.init()
	return s
}

func (s *Standard) init() {
	s.notes = make(map[int]map[string]Note)
	s.reprs = make(map[string]Note)
	s.indexes = make(map[int]Note)
	afreq := 55.0
	num := 1
	for oct := 1; oct <= 7; oct++ {
		// s.notes[oct] = s.buildOctaveScale(afreq)
		num = s.buildOctaveRepr(oct, afreq, num)
		afreq *= 2.0
	}
}

func (s *Standard) buildOctaveRepr(oct int, afreq float64, num int) int {
	s.notes[oct] = make(map[string]Note)
	var idx int

	a := afreq
	idx = num + 9
	nt := Note{idx, a}
	s.reprs[fmt.Sprintf("A%d", oct)] = nt
	s.notes[oct]["A"] = nt
	s.indexes[idx] = nt

	as := a * HalfStep
	idx = num + 10
	nt = Note{idx, as}
	s.reprs[fmt.Sprintf("As%d", oct)] = nt
	s.reprs[fmt.Sprintf("Bb%d", oct)] = nt
	s.notes[oct]["a"] = nt
	s.indexes[idx] = nt

	b := as * HalfStep
	idx = num + 11
	nt = Note{idx, b}
	s.reprs[fmt.Sprintf("B%d", oct)] = nt
	s.notes[oct]["B"] = nt
	s.indexes[idx] = nt

	ab := a / HalfStep
	idx = num + 8
	nt = Note{idx, ab}
	s.reprs[fmt.Sprintf("Ab%d", oct)] = nt
	s.reprs[fmt.Sprintf("Gs%d", oct)] = nt
	s.notes[oct]["g"] = nt
	s.indexes[idx] = nt

	g := ab / HalfStep
	idx = num + 7
	nt = Note{idx, g}
	s.reprs[fmt.Sprintf("G%d", oct)] = nt
	s.notes[oct]["G"] = nt
	s.indexes[idx] = nt

	gb := g / HalfStep
	idx = num + 6
	nt = Note{idx, gb}
	s.reprs[fmt.Sprintf("Gb%d", oct)] = nt
	s.reprs[fmt.Sprintf("Fs%d", oct)] = nt
	s.notes[oct]["f"] = nt
	s.indexes[idx] = nt

	f := gb / HalfStep
	idx = num + 5
	nt = Note{idx, f}
	s.reprs[fmt.Sprintf("F%d", oct)] = nt
	s.notes[oct]["F"] = nt
	s.indexes[idx] = nt

	e := f / HalfStep
	idx = num + 4
	nt = Note{idx, e}
	s.reprs[fmt.Sprintf("E%d", oct)] = nt
	s.notes[oct]["E"] = nt
	s.indexes[idx] = nt

	eb := e / HalfStep
	idx = num + 3
	nt = Note{idx, eb}
	s.reprs[fmt.Sprintf("Eb%d", oct)] = nt
	s.reprs[fmt.Sprintf("Ds%d", oct)] = nt
	s.notes[oct]["d"] = nt
	s.indexes[idx] = nt

	d := eb / HalfStep
	idx = num + 2
	nt = Note{idx, d}
	s.reprs[fmt.Sprintf("D%d", oct)] = nt
	s.notes[oct]["D"] = nt
	s.indexes[idx] = nt

	db := d / HalfStep
	idx = num + 1
	nt = Note{idx, db}
	s.reprs[fmt.Sprintf("Db%d", oct)] = nt
	s.reprs[fmt.Sprintf("Cs%d", oct)] = nt
	s.notes[oct]["c"] = nt
	s.indexes[idx] = nt

	c := db / HalfStep
	idx = num
	nt = Note{idx, c}
	s.reprs[fmt.Sprintf("C%d", oct)] = nt
	s.notes[oct]["C"] = nt
	s.indexes[idx] = nt

	return num + 12
}

func (s *Standard) Note(octave int, note string) (Note, bool) {
	hz, ok := s.notes[octave][note]
	return Note(hz), ok
}

func (s *Standard) Parse(str string) (Note, bool) {
	n, ok := s.reprs[str]
	return n, ok
}

func (e *Standard) ByNumber(n int) (Note, bool) {
	nt, ok := e.indexes[n]
	return nt, ok
}

func (s *Standard) Edo() int {
	return 12
}

func (s *Standard) Std() []byte {
	return std.StdEdo12
}
