package filters

import "fmt"

func float64Of(x any) float64 {
	switch a := x.(type) {
	case float64:
		return a
	case int:
		return float64(a)
	default:
		panic(fmt.Errorf("Cannot convert to float64: %v (%T)", x, x))
	}
}
