package notes

// 12th root of 2.0
const HalfStep = 1.0594630943593

// Standard
type Standard struct {
	notes map[int]map[string]float64
}

func NewStandard() *Standard {
	s := &Standard{}
	s.init()
	return s
}

func (s *Standard) init() {
	s.notes = make(map[int]map[string]float64)
	s.notes[1] = s.buildOctaveScale(55.0)
	s.notes[2] = s.buildOctaveScale(110.0)
	s.notes[3] = s.buildOctaveScale(220.0)
	s.notes[4] = s.buildOctaveScale(440.0)
	s.notes[5] = s.buildOctaveScale(880.0)
	s.notes[6] = s.buildOctaveScale(880.0 * 2.0)
	s.notes[7] = s.buildOctaveScale(880.0 * 4.0)
}

func (s *Standard) buildOctaveScale(afreq float64) map[string]float64 {
	m := map[string]float64{
		"A": afreq,
	}

	as := afreq * HalfStep
	m["a"] = as
	b := as * HalfStep
	m["B"] = b

	gs := afreq / HalfStep
	m["g"] = gs
	g := gs / HalfStep
	m["G"] = g
	fs := g / HalfStep
	m["f"] = fs
	f := fs / HalfStep
	m["F"] = f
	e := f / HalfStep
	m["E"] = e

	ds := e / HalfStep
	m["d"] = ds
	d := ds / HalfStep
	m["D"] = d
	cs := d / HalfStep
	m["c"] = cs
	c := cs / HalfStep
	m["C"] = c

	return m
}

func (s *Standard) Note(octave int, note string) (hz float64, ok bool) {
	hz, ok = s.notes[octave][note]
	return
}
