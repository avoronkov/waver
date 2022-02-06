package notes

// 12th root of 2.0
const HalfStep = 1.0594630943593

// Standard
type Standard struct {
	notes map[int]map[string]float64
}

var _ EdoScale = (*Standard)(nil)

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

// const HalfStep = 1.0594630943593
func (s *Standard) buildOctaveScale(afreq float64) map[string]float64 {
	m := map[string]float64{"A": afreq}

	m["a"] = m["A"] * HalfStep
	m["B"] = m["a"] * HalfStep

	m["g"] = m["A"] / HalfStep
	m["G"] = m["g"] / HalfStep
	m["f"] = m["G"] / HalfStep
	m["F"] = m["f"] / HalfStep
	m["E"] = m["F"] / HalfStep

	m["d"] = m["E"] / HalfStep
	m["D"] = m["d"] / HalfStep
	m["c"] = m["D"] / HalfStep
	m["C"] = m["c"] / HalfStep

	return m
}

func (s *Standard) Note(octave int, note string) (hz float64, ok bool) {
	hz, ok = s.notes[octave][note]
	return
}

func (s *Standard) Edo() int {
	return 12
}
