package common

import "gitlab.com/avoronkov/waver/lib/seq/types"

type Num int64

var _ types.Value = Num(0)

func (n Num) ToInt64List() []int64 {
	return []int64{int64(n)}
}

type List []Num

var _ types.Value = List(nil)

func (l List) ToInt64List() (res []int64) {
	for _, n := range l {
		res = append(res, int64(n))
	}
	return
}
