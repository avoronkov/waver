package main

type Value interface {
	ToInt64List() []int64
}

type Num int64

func (n Num) ToInt64List() []int64 {
	return []int64{int64(n)}
}

type List []Num

func (l List) ToInt64List() (res []int64) {
	for _, n := range l {
		res = append(res, int64(n))
	}
	return
}
