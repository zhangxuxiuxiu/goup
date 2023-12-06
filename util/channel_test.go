package util

import (
	"fmt"
	"testing"
	"time"
)

func cprint(arr <-chan int) {
	fmt.Print("[")
	for v := range arr {
		fmt.Printf("%d,", v)
	}
	fmt.Print("]\n")
}

func cprint2(arr <-chan chan int) {
	fmt.Print("[\n")
	for v := range arr {
		cprint(v)
	}
	fmt.Print("]\n\n")
}

func TestCBatch(x *testing.T) {
	cprint2(CBatch(ToChan([]int{}), 1))
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8}
	cprint2(CBatch(ToChan(arr), 1))
	cprint2(CBatch(ToChan(arr), 2))
	cprint2(CBatch(ToChan(arr), 3))
	cprint2(CBatch(ToChan(arr), 4))
	cprint2(CBatch(ToChan(arr), 5))
	cprint2(CBatch(ToChan(arr), 6))
}

func TestCUnique(x *testing.T) {
	arr1 := []int{}
	cprint(CUnique(ToChan(arr1)))
	arr2 := []int{1}
	cprint(CUnique(ToChan(arr2)))
	arr3 := []int{1, 1}
	cprint(CUnique(ToChan(arr3)))
	arr4 := []int{1, 2, 3}
	cprint(CUnique(ToChan(arr4)))
	arr5 := []int{1, 2, 2, 3}
	cprint(CUnique(ToChan(arr5)))
}

func TestCEqual(x *testing.T) {
	arr1 := []int{}
	arr2 := []int{1}
	arr3 := []int{1, 1}
	arr4 := []int{1, 2, 3}
	arr5 := []int{1, 2, 2, 3}
	fmt.Println(CEqual(ToChan(arr1), ToChan(arr1), nil, nil))
	fmt.Println(!CEqual(ToChan(arr1), ToChan(arr2), nil, nil))
	fmt.Println(CEqual(ToChan(arr2), ToChan(arr2), nil, nil))
	fmt.Println(CEqual(ToChan(arr3), ToChan(arr3), nil, nil))
	fmt.Println(!CEqual(ToChan(arr4), ToChan(arr5), nil, nil))
}

func TestCRange(x *testing.T) {
	for x := range CRange() {
		fmt.Printf("%d,", x)
	}
	fmt.Println()

	for x := range CRange(10) {
		fmt.Printf("%d,", x)
	}
	fmt.Println()

	for x := range CRange(10, 0) {
		fmt.Printf("%d,", x)
	}
	fmt.Println()

	for x := range CRange(10, 0, 2) {
		fmt.Printf("%d,", x)
	}
	fmt.Println()

	for x := range CRange(10, 0, 2, 1) {
		fmt.Printf("%d,", x)
	}
	fmt.Println()
}

func sleep(msec int64) func() int64 {
	return func() int64 {
		time.Sleep(time.Millisecond * time.Duration(msec))
		return msec
	}
}

func less(upper int64) func(int64) bool {
	return func(v int64) bool {
		return v < upper
	}
}

func TestWhenAny(x *testing.T) {
	less3 := less(3)
	less10 := less(10)
	less15 := less(10)

	empty := []func() int64{}
	fmt.Printf("empty:%v\n", <-WhenAny(empty, less10))

	rpc1 := []func() int64{sleep(5)}
	fmt.Printf("5<10:%d\n", <-WhenAny(rpc1, less10))
	fmt.Printf("5<3:%v\n", <-WhenAny(rpc1, less3))

	rpc2 := []func() int64{sleep(5), sleep(11)}
	fmt.Printf("5,11<3:%v\n", <-WhenAny(rpc2, less3))
	fmt.Printf("5,11<10:%d\n", <-WhenAny(rpc2, less10))
	fmt.Printf("5,11<15:%d\n", <-WhenAny(rpc2, less15))
}
