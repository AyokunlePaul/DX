package utils

func Contains[T comparable](data []T, value T) bool {
	for _, a := range data {
		if a == value {
			return true
		}
	}
	return false
}
