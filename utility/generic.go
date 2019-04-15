package utility

//Contains checks if element is in container
func Contains(container []int, element int) bool {
	for _, a := range container {
		if a == element {
			return true
		}
	}
	return false
}

//Min returns minimum from a and b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
