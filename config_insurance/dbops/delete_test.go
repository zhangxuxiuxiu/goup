package dbops

import (
	"fmt"
	"testing"
)

func TestDelete(t *testing.T) {
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
	fmt.Printf("gen soft delete for %v\n\ngenerated sql is:\n%s\n", p1, Delete(&p1, "product_tab"))
}
