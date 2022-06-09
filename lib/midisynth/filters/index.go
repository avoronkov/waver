package filters

type FilterCreator interface {
	Create(options any) (Filter, error)
}

var Filters = map[string]FilterCreator{
	"adsr":       AdsrFilter{},
	"delay":      DelayFilter{},
	"distortion": DistortionFilter{},
	"vibrato":    VibratoFilter{},
	"am":         Ring{},
	"timeshift":  TimeShift{},
	"harmonizer": Harmonizer{},
	"flanger":    Flanger{},
	"exp":        Exponent{},
	"movexp":     MovingExponent{},
	"ratio":      Ratio{},
}
