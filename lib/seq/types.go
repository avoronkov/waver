package seq

import "gitlab.com/avoronkov/waver/lib/seq/types"

type assignment struct {
	name    string
	valueFn types.ValueFn
}

type assignments []assignment

func (a assignments) Get(name string) (types.ValueFn, bool) {
	for _, as := range a {
		if as.name == name {
			return as.valueFn, true
		}
	}
	return nil, false
}
