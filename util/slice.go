package util

import (
	"math/rand"
	"strings"
)

func Any[T any](arr []T, pred func(T) bool) bool {
	for _, v := range arr {
		if pred(v) {
			return true
		}
	}
	return false
}

func All[T any](arr []T, pred func(T) bool) bool {
	return !Any(arr, Not(pred))
}

func Equal[T comparable](src, dest []T) bool {
	if len(src) != len(dest) {
		return false
	}
	for idx, v := range src {
		if v != dest[idx] {
			return false
		}
	}
	return true
}

func Find[T any](arr []T, pred func(T) bool) int {
	for idx, e := range arr {
		if pred(e) {
			return idx
		}
	}
	return -1
}

func FindV[T comparable](arr []T, v T) int {
	return Find(arr, EqualTo(v))
}

func Remove[T any](arr []T, pos int) (bool, []T) {
	if pos < 0 || pos >= len(arr) {
		return false, arr
	}
	arr[pos] = arr[len(arr)-1]
	return true, arr[:len(arr)-1]
}

func Insert[T any](arr []T, pos int, val T) []T {
	var tmp T
	arr = append(arr, tmp)
	copy(arr[pos+1:], arr[pos:])
	arr[pos] = val
	return arr
}

func Reverse[T any](arr []T) []T {
	for i := len(arr)/2 - 1; i >= 0; i-- {
		opp := len(arr) - 1 - i
		arr[i], arr[opp] = arr[opp], arr[i]
	}
	return arr
}

func Shuffle[T any](arr []T) []T {
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func Batch[T any](arr []T, batchSize int) [][]T {
	batches := make([][]T, 0, (len(arr)+batchSize-1)/batchSize)

	for batchSize < len(arr) {
		arr, batches = arr[batchSize:], append(batches, arr[0:batchSize])
	}
	if len(arr) > 0 {
		batches = append(batches, arr)
	}
	return batches
}

// Unique  arr must be sorted
func Unique[T comparable](arr []T) []T {
	if len(arr) == 0 {
		return arr
	}
	j := 1
	for i := 1; i < len(arr); i++ {
		if arr[i-1] == arr[i] {
			continue
		}
		arr[j] = arr[i]
		j++
	}
	return arr[:j]
}

// Filter, Map, Reduce
func Filter[T any](arr []T, pred func(T) bool) []T {
	n := 0
	for _, v := range arr {
		if pred(v) {
			arr[n] = v
			n++
		}
	}
	return arr[:n]
}

func Map[T, U any](src []T, fn func(T) U) []U {
	dest := make([]U, 0, len(src))
	for _, v := range src {
		dest = append(dest, fn(v))
	}
	return dest
}

func FlatMap[T, U any](src []T, fn func(T) []U) []U {
	dest := make([]U, 0)
	for _, v := range src {
		if us := fn(v); len(us) != 0 {
			dest = append(dest, us...)
		}
	}
	return dest
}

func Accumulate[T, U any](arr []T, fn func(U, T) U, base U) U {
	for _, v := range arr {
		base = fn(base, v)
	}
	return base
}

func Reduce[T any](arr []T, fn func(T, T) T) T {
	var zero T
	return Accumulate(arr, fn, zero)
}

func Copy[T any](src []T) []T {
	dest := make([]T, len(src))
	copy(dest, src)
	return dest
}

func ToSlice[T any](c chan T) []T {
	s := make([]T, 0)
	for v := range c {
		s = append(s, v)
	}
	return s
}

type Stringer interface {
	comparable
	String() string
}

func Join[T Stringer](arr []T, sep string) string {
	var builder strings.Builder

	for _, r := range arr {
		if builder.Len() != 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(r.String())
	}
	return builder.String()
}
