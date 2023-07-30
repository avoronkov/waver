package utils

func SliceMap[T any, R any](slice []T, fn func(x T) R) []R {
	res := make([]R, len(slice))
	for i, v := range slice {
		res[i] = fn(v)
	}
	return res
}
