package hlp

func PtrTo[T any](value T) *T {
	return &value
}
