package dbops

import (
	"fmt"
	"testing"
)

type X struct {
	Id   int
	Name string
}

func (a *X) Equal(b *X) bool {
	fmt.Printf("compare in X.Equal\n")
	return *a == *b
}

type Y struct {
	Id   int `equal:"-"`
	Name string
}

type Z struct {
	Id   int
	Name string `equal:"-"`
}

func TestEqual(*testing.T) {
	a, b, c := X{Id: 1, Name: ""}, X{Id: 3, Name: ""}, X{Id: 3, Name: ""}
	if Equal(&a, &b) {
		fmt.Printf("%d!=%d\n", a.Id, b.Id)
	}
	if !Equal(&c, &b) {
		fmt.Printf("%d!=%d\n", c.Id, b.Id)
	}

	y1, y2, y3, y4 := Y{Id: 1, Name: "y"}, Y{Id: 1, Name: "y"}, Y{Id: 2, Name: "y"}, Y{Id: 2, Name: "y2"}
	if !Equal(y1, y2) || !Equal(y2, y3) || Equal(y3, y4) {
		fmt.Printf("expect y1==y2==y3!=y4\n")
	}

	z1, z2, z3, z4 := Z{Id: 1, Name: "y"}, Z{Id: 1, Name: "y"}, Z{Id: 1, Name: "y2"}, Z{Id: 2, Name: "y2"}
	if !Equal(z1, z2) || !Equal(z2, z3) || Equal(z3, z4) {
		fmt.Printf("expect z1==z2==z3!=z4\n")
	}
}
