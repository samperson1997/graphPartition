package pshp


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
