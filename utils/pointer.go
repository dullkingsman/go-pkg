package utils

func PtrOf[T any](value T) *T {
	return &value
}

func ValueOr[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}
