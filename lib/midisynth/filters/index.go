package filters

type FilterCreator interface {
	Create(options any) (Filter, error)
}

var FilterCreators = map[string]FilterCreator{
	"adsr":       AdsrFilter{},
	"delay":      DelayFilter{},
	"distortion": DistortionFilter{},
	"dist":       DistortionFilter{},
	"vibrato":    VibratoFilter{},
	"am":         Ring{},
	"timeshift":  TimeShift{},
	"harmonizer": Harmonizer{},
	"harm":       Harmonizer{},
	"flanger":    Flanger{},
	"exp":        Exponent{},
	"movexp":     MovingExponent{},
	"pan":        Pan{},
	"movpan":     MovingPan{},
	"swingexp":   SwingExp{},
}

type NewFilter interface {
	New() Filter
}

var Filters = map[string]NewFilter{
	"ratio": Ratio{},
}
