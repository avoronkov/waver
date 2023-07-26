package surfer

func SlicesFromSamples(data []float64) [][]float64 {
	result := [][]float64{}
	idx := 0
	for {
		var slice []float64
		idx, slice = getSlice(data, idx)
		if len(slice) == 0 {
			break
		}
		result = append(result, slice)
	}
	return result
}

func getSlice(data []float64, idx int) (int, []float64) {
	res := []float64{}
	l := len(data)
	// handle zero samples
	for idx < l && data[idx] == 0.0 {
		res = append(res, 0.0)
		idx++
	}

	if idx >= l {
		return idx, nil
	}

	sign := 1.0
	if data[idx] < 0 {
		sign = -1.0
	}
	idx, a := getSignedSlice(data, idx, sign)
	if len(a) == 0 {
		return idx, res
	}
	res = append(res, a...)

	// handle zero samples
	for idx < l && data[idx] == 0.0 {
		res = append(res, 0.0)
		idx++
	}

	idx, b := getSignedSlice(data, idx, -1.0*sign)
	if len(b) == 0 {
		return idx, res
	}
	res = append(res, b...)
	return idx, res
}

func getSignedSlice(data []float64, idx int, sign float64) (int, []float64) {
	res := []float64{}
	l := len(data)
	for idx < l && data[idx]*sign > 0 {
		res = append(res, data[idx])
		idx++
	}
	return idx, res
}
