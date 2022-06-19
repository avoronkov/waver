package filters

type FilterCreator interface {
	Create(options any) (Filter, error)
}

var FilterCreators = map[string]FilterCreator{
	"harmonizer": Harmonizer{},
	"harm":       Harmonizer{},
	"movexp":     MovingExponent{},
}

type NewFilter interface {
	New() Filter
}

var Filters = map[string]NewFilter{
	"adsr":       AdsrFilter{},
	"am":         Ring{},
	"delay":      DelayFilter{},
	"distortion": Distortion{},
	"dist":       Distortion{},
	"exp":        Exponent{},
	"flanger":    Flanger{},
	"pan":        Pan{},
	"ratio":      Ratio{},
	"swingexp":   SwingExp{},
	"swingpan":   SwingPan{},
	"timeshift":  TimeShift{},
	"vibrato":    VibratoFilter{},
}
