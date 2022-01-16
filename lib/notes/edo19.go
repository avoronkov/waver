package notes

// 19th root of 2.0
const Edo19Step = 1.0371550444462

type Edo19 struct {
	notes map[int]map[string]float64
}

var _ (EdoScale) = (*Edo19)(nil)

func NewEdo19() *Edo19 {
	e := &Edo19{}
	e.init()
	return e
}

func (e *Edo19) init() {
	e.notes = make(map[int]map[string]float64)
	e.notes[1] = e.buildOctaveScale(55.0)
	e.notes[2] = e.buildOctaveScale(110.0)
	e.notes[3] = e.buildOctaveScale(220.0)
	e.notes[4] = e.buildOctaveScale(440.0)
	e.notes[5] = e.buildOctaveScale(880.0)
	e.notes[6] = e.buildOctaveScale(880.0 * 2.0)
	e.notes[7] = e.buildOctaveScale(880.0 * 4.0)
}

func (*Edo19) buildOctaveScale(afreq float64) map[string]float64 {
	m := map[string]float64{
		"A": afreq,
	}

	as := afreq * Edo19Step
	m["a"] = as
	bf := as * Edo19Step
	m["9"] = bf
	b := bf * Edo19Step
	m["B"] = b
	bs := b * Edo19Step
	m["b"] = bs

	af := afreq / Edo19Step
	m["8"] = af
	gs := af / Edo19Step
	m["g"] = gs
	g := gs / Edo19Step
	m["G"] = g
	gf := g / Edo19Step
	m["7"] = gf
	fs := gf / Edo19Step
	m["f"] = fs
	f := fs / Edo19Step
	m["F"] = f

	es := f / Edo19Step
	m["e"] = es
	e := es / Edo19Step
	m["E"] = e
	ef := e / Edo19Step
	m["6"] = ef
	ds := ef / Edo19Step
	m["d"] = ds
	d := ds / Edo19Step
	m["D"] = d
	df := d / Edo19Step
	m["5"] = df
	cs := df / Edo19Step
	m["c"] = cs
	c := cs / Edo19Step
	m["C"] = c

	return m
}

func (e *Edo19) Note(octave int, note string) (hz float64, ok bool) {
	hz, ok = e.notes[octave][note]
	return
}

func (e *Edo19) Edo() int {
	return 19
}
