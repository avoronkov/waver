package config

import (
	"github.com/avoronkov/waver/lib/midisynth/filters"
	"github.com/avoronkov/waver/lib/midisynth/waves"
)

type FilterCreator interface {
	Create(options any) (filters.Filter, error)
}

var filterCreators = map[string]FilterCreator{
	"adsr":       filters.AdsrFilter{},
	"delay":      filters.DelayFilter{},
	"distortion": filters.DistortionFilter{},
	"vibrato":    filters.VibratoFilter{},
	"am":         filters.Ring{},
	"timeshift":  filters.TimeShift{},
	"harmonizer": filters.Harmonizer{},
	"flanger":    filters.Flanger{},
	"exp":        filters.Exponent{},
	"movexp":     filters.MovingExponent{},
	"ratio":      filters.Ratio{},
}

var waveFunctions = map[string]waves.Wave{
	"sine":     waves.Sine,
	"square":   waves.Square,
	"triangle": waves.Triangle,
	"saw":      waves.Saw,
	"semisine": waves.SemiSine,
}
