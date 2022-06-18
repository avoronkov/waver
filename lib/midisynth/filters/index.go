package filters

type FilterCreator interface {
	Create(options any) (Filter, error)
}

var FilterCreators = map[string]FilterCreator{
	"distortion": DistortionFilter{},
	"dist":       DistortionFilter{},
	"vibrato":    VibratoFilter{},
	"am":         Ring{},
	"timeshift":  TimeShift{},
	"harmonizer": Harmonizer{},
	"harm":       Harmonizer{},
	"flanger":    Flanger{},
	"movexp":     MovingExponent{},
	"pan":        Pan{},
	"movpan":     MovingPan{},
	"swingexp":   SwingExp{},
}

type NewFilter interface {
	New() Filter
}

var Filters = map[string]NewFilter{
	"adsr":  AdsrFilter{},
	"delay": DelayFilter{},
	"exp":   Exponent{},
	"ratio": Ratio{},
}
