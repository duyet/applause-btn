package utils

// Clamp get next number in range
func Clamp(value, lower, upper int) int {
	return max(lower, min(value, upper))
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
