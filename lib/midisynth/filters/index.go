package filters

type FilterCreator interface {
	Create(options any) (Filter, error)
}

type NewFilter interface {
	New() Filter
}

var Filters = map[string]NewFilter{
	"8bit":       EightBit{},
	"adsr":       AdsrFilter{},
	"am":         Ring{},
	"delay":      DelayFilter{},
	"dist":       Distortion{},
	"distortion": Distortion{},
	"exp":        Exponent{},
	"flanger":    Flanger{},
	"harmonizer": Harmonizer{},
	"harm":       Harmonizer{},
	"movexp":     MovingExponent{},
	"pan":        Pan{},
	"ratio":      Ratio{},
	"swingexp":   SwingExp{},
	"swingpan":   SwingPan{},
	"timeshift":  TimeShift{},
	"vibrato":    VibratoFilter{},
	"code":       Custom{},
}
