package main

import (
	"log"
	"math"
)

type SamplesStorage struct {
	buf []float64
}

func (s *SamplesStorage) AddSample(n int, value float64) {
	if n >= len(s.buf) {
		newlen := (n + 1) * 2
		newbuf := make([]float64, newlen)
		copy(newbuf, s.buf)
		s.buf = newbuf
	}

	s.buf[n] += value
}

func (s *SamplesStorage) Normalize() {
	l, r := s.samplesRange()
	max := s.findMaxValue(l, r)
	if max == 0 {
		return
	}
	for i := l; i < r; i++ {
		s.buf[i] = s.buf[i] / max
	}
}

func (s *SamplesStorage) findMaxValue(l, r int) float64 {
	max := 0.0
	for i := l; i < r; i++ {
		if mv := math.Abs(s.buf[i]); mv > max {
			max = mv
		}
	}
	return max
}

// [min, max)
func (s *SamplesStorage) samplesRange() (min int, max int) {
	l := len(s.buf)
	for ; min < l; min++ {
		if s.buf[min] != 0.0 {
			break
		}
	}
	if min == l {
		log.Printf("Empy samples storage: non-zero samples not found (%v)", l)
		return 0, 0
	}

	for max = l - 1; max >= 0; max-- {
		if s.buf[max] != 0.0 {
			break
		}
	}
	max++
	return min, max
}

const maxInt16Amp = (1 << 15) - 1

func (s *SamplesStorage) ToInt16List() []int16 {
	l, r := s.samplesRange()
	size := r - l
	if size <= 0 {
		return nil
	}
	res := make([]int16, size)
	j := 0
	for i := l; i < r; i++ {
		res[j] = int16(s.buf[i] * maxInt16Amp)
		j++
	}
	return res
}
