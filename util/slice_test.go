package util

import (
	"fmt"
	"testing"
	"unsafe"
)

func print(arr []int) {
	fmt.Print("[")
	for _, v := range arr {
		fmt.Printf("%d,", v)
	}
	fmt.Print("]\n")
}

func print2(arr [][]int) {
	fmt.Print("[\n")
	for _, v := range arr {
		print(v)
	}
	fmt.Print("]\n\n")
}

func TestUnique(x *testing.T) {
	arr1 := []int{}
	print(Unique(arr1))
	arr2 := []int{1}
	print(Unique(arr2))
	arr3 := []int{1, 1}
	print(Unique(arr3))
	arr4 := []int{1, 2, 3}
	print(Unique(arr4))
}

func TestBatch(x *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7}
	print2(Batch(arr, 1))
	print2(Batch(arr, 2))
	print2(Batch(arr, 3))
	print2(Batch(arr, 4))
	print2(Batch(arr, 5))
	print2(Batch(arr, 6))
}

type Outer struct {
	a int
	b float64
}

type Outer2 struct {
	ot Outer
	a  int
	b  float64
}

func TestStructFieldNoCopy(t *testing.T) {
	out := Outer2{Outer{2, 0.9}, 1, 0.2}
	if &out.ot.b != &out.ot.b {
		t.Errorf("copy happened")
	} else {
		t.Logf("no copy with addr:%d,%d", unsafe.Pointer(&out.ot.b), unsafe.Pointer(&out.ot.b))
	}

	key2conf := make(map[string]Outer2)
	key2conf["id"] = out
	out1 := key2conf["id"]
	out2 := key2conf["id"]

	if &out1 != &out2 {
		t.Logf("copy happened with addr:%d,%d", unsafe.Pointer(&out1), unsafe.Pointer(&out2))
	} else {
		t.Errorf("no copy with addr:%d,%d", unsafe.Pointer(&out1), unsafe.Pointer(&out2))
	}

	//if &(key2conf["id"]) != &(key2conf["id"]) { //Cannot take the address of 'key2conf["id"]'
	//	t.Logf("copy happened with addr:%d,%d", unsafe.Pointer(&out1), unsafe.Pointer(&out2))
	//} else {
	//	t.Errorf("no copy with addr:%d,%d", unsafe.Pointer(&out1), unsafe.Pointer(&out2))
	//}
}
