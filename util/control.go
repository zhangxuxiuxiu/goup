package util

func IfElse[T any](b bool, a1, a2 T) T {
	if b {
		return a1
	} else {
		return a2
	}
}
