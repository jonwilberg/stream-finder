package datatools

// Return unique items from a slice, maintaining order.
func Unique[T comparable](input []T) []T {
	seen := make(map[T]struct{}, len(input))
	result := make([]T, 0, len(input))
	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
