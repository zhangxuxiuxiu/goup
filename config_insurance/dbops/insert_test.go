package dbops

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	var p1 = Product{Pid: 123,
		Name: "pro1",
		Plan: Plan{
			Pid:  345,
			Name: "plan1",
		},
		Benefits: []Benefit{
			{Bid: 1, Name: "benefit1"},
			{Bid: 2, Name: "benefit2"},
		},
	}
	fmt.Printf("gen insert for %v\n\ngenerated sql is:\n%s\n", p1, Insert(&p1, "product_tab"))
}
