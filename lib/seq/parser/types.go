package parser

import "gitlab.com/avoronkov/waver/lib/seq/types"

type Seq interface {
	Add(types.Signaler)
	Commit() error
}
