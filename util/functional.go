package util

// functional Not, Or, And
func Not[T any](fn func(T) bool) func(T) bool {
	return func(v T) bool {
		return !fn(v)
	}
}

func Not2[T, U any](fn func(T, U) bool) func(T, U) bool {
	return func(t T, u U) bool {
		return !fn(t, u)
	}
}

func Or[T any](functors []func(T) bool) func(T) bool {
	return func(v T) bool {
		return Any(functors, func(fn func(T) bool) bool {
			return fn(v)
		})
	}
}

func Or2[T, U any](functors []func(T, U) bool) func(T, U) bool {
	return func(v T, u U) bool {
		return Any(functors, func(fn func(T, U) bool) bool {
			return fn(v, u)
		})
	}
}

func And[T any](functors []func(T) bool) func(T) bool {
	return func(v T) bool {
		return All(functors, func(fn func(T) bool) bool {
			return fn(v)
		})
	}
}

func And2[T, U any](functors []func(T, U) bool) func(T, U) bool {
	return func(v T, u U) bool {
		return All(functors, func(fn func(T, U) bool) bool {
			return fn(v, u)
		})
	}
}

func IsErr(e error) bool {
	return e != nil
}

func NilErr(e error) bool {
	return e == nil
}

func True[T any](v T) bool {
	return true
}

func False[T any](v T) bool {
	return false
}

func Equaler[T comparable](a, b T) bool {
	return a == b
}

func Unequaler[T comparable](a, b T) bool {
	return a != b
}

func Identity[T any](v T) T {
	return v
}

func EqualTo[T comparable](a T) func(T) bool {
	return func(b T) bool {
		return a == b
	}
}

func Value[T any](v T) func() T {
	return func() T {
		return v
	}
}
