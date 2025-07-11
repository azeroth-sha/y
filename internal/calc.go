package internal

// Or returns y if cond is true, n otherwise.
func Or[T any](cond bool, y, n T) T {
	if cond {
		return y
	}
	return n
}
