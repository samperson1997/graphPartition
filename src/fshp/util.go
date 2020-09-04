package fshp

func imin(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
func imax(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
