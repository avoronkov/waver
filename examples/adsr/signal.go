package main

import "math"

type Signal struct {
	AttackLevel int16
	AttackLen   int
	DecayLevel  int16
	DecayLen    int
	SusteinLen  int
	ReleaseLen  int
}

type Sampler interface {
	AddStereoSample(value int16)
}

func (s *Signal) PutSignal(hz float64, sampler Sampler) {
	samplesPerSecond := 44100
	var waveDuration float64 = float64(samplesPerSecond) / hz
	// attack
	i := 0
	for ; i < int(s.AttackLen); i++ {
		x := 2.0 * math.Pi * float64(i) / waveDuration
		amp := float64(i) * float64(s.AttackLevel) / float64(s.AttackLen-1)
		l := int16(amp * math.Sin(x))
		sampler.AddStereoSample(l)
	}

	// decay
	for j := 0; j < int(s.DecayLen); j++ {
		i++
		x := 2.0 * math.Pi * float64(i) / waveDuration
		amp := float64(s.AttackLevel) - float64(s.AttackLevel-s.DecayLevel)*float64(j)/float64(s.DecayLen)
		l := int16(amp * math.Sin(x))
		sampler.AddStereoSample(l)
	}

	// sustain
	for j := 0; j < int(s.SusteinLen); j++ {
		i++
		x := 2.0 * math.Pi * float64(i) / waveDuration
		l := int16(float64(s.DecayLevel) * math.Sin(x))
		sampler.AddStereoSample(l)
	}

	// release
	for j := 0; j < int(s.ReleaseLen); j++ {
		i++
		x := 2.0 * math.Pi * float64(i) / waveDuration
		amp := float64(s.DecayLevel) - float64(s.DecayLevel)*float64(j)/float64(s.ReleaseLen)
		l := int16(amp * math.Sin(x))
		sampler.AddStereoSample(l)
	}
}
