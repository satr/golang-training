package loopsandrange

// Sum returns the total of all values.
func Sum(values []int) int {
	total := 0
	for _, value := range values {
		total += value
	}
	return total
}

// RuneCounts counts each Unicode code point in text.
func RuneCounts(text string) map[rune]int {
	counts := make(map[rune]int)
	for _, value := range text {
		counts[value]++
	}
	return counts
}
