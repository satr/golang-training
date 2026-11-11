package slicesandarrays

// ReverseFive returns a reversed copy of a five-element array.
func ReverseFive(values [5]int) [5]int {
	var reversed [5]int
	for index, value := range values {
		reversed[len(values)-1-index] = value
	}
	return reversed
}

// AppendCopy combines slices without sharing writable storage with either input.
func AppendCopy(values, extra []int) []int {
	if values == nil && extra == nil {
		return nil
	}
	combined := make([]int, 0, len(values)+len(extra))
	combined = append(combined, values...)
	combined = append(combined, extra...)
	return combined
}
