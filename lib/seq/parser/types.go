package parser

import "gitlab.com/avoronkov/waver/lib/seq/types"

type Seq interface {
	Add(types.Signaler)
	Commit() error
	Assign(name string, value types.ValueFn)
}
