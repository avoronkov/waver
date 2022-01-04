package waves

type Wave interface {
	// Return sample volume at moment t (in seconds) in range [-1.0, 1.0]
	Value(t float64) float64
}
