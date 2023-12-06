package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDeepEqualMap(*testing.T) {
	var v int = 0
	m1 := map[string]*int{
		"k1": nil,
		"k2": &v,
	}
	m2 := map[string]*int{
		"k1": nil,
		"k2": &v,
	}
	fmt.Printf("equal:%v\n", reflect.DeepEqual(m1, m2))
}
