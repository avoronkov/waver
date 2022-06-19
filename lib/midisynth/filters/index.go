package filters

type FilterCreator interface {
	Create(options any) (Filter, error)
}

var FilterCreators = map[string]FilterCreator{
	"vibrato":    VibratoFilter{},
	"am":         Ring{},
	"timeshift":  TimeShift{},
	"harmonizer": Harmonizer{},
	"harm":       Harmonizer{},
	"movexp":     MovingExponent{},
	"swingexp":   SwingExp{},
}

type NewFilter interface {
	New() Filter
}

var Filters = map[string]NewFilter{
	"adsr":       AdsrFilter{},
	"delay":      DelayFilter{},
	"distortion": Distortion{},
	"dist":       Distortion{},
	"exp":        Exponent{},
	"flanger":    Flanger{},
	"pan":        Pan{},
	"ratio":      Ratio{},
	"swingpan":   SwingPan{},
}
