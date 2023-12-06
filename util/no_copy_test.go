package util

import (
	"fmt"
	"testing"
)

type Tester struct {
	NoCopy
}

func TestNoCopy(t *testing.T) {
	var nc Tester
	nc2 := nc
	nc = nc2
	fmt.Printf("%T", nc2)
}
