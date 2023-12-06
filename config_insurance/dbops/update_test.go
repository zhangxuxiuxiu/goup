package dbops

import (
	"fmt"
	"testing"
)

type Product struct {
	Pid      uint64 `id:""`
	Name     string
	Plan     Plan
	Benefits []Benefit
}

func (p *Product) Id() uint64 {
	return p.Pid
}

//func (p *Product) GenUpdateSql() string {
//	return fmt.Sprintf("update product_tab set name='%s' where id=%d;", p.Name, p.Pid)
//}
//
//func (p *Product) GenDeleteSql() string {
//	return fmt.Sprintf("update product_tab set is_deleted=1 where id=%d;", p.Pid)
//}
//
//func (p *Product) GenInsertSql() string {
//	return fmt.Sprintf("insert into product_tab values(%d,'%s');", p.Pid, p.Name)
//}

type Plan struct {
	Pid  uint64 `id:""`
	Name string
}

//func (p *Plan) Id() uint64 {
//	return p.Pid
//}
//
//func (p Plan) GenUpdateSql() string {
//	return fmt.Sprintf("update plan_tab set name='%s' where id=%d;", p.Name, p.Pid)
//}
//
//func (p Plan) GenDeleteSql() string {
//	return fmt.Sprintf("update plan_tab set is_deleted=1 where id=%d;", p.Pid)
//}
//
//func (p Plan) GenInsertSql() string {
//	return fmt.Sprintf("insert into plan_tab values(%d,'%s');", p.Pid, p.Name)
//}

type Benefit struct {
	Bid  uint64 `id:""`
	Name string
}

//func (p Benefit) Id() uint64 {
//	return p.Bid
//}
//
//func (p Benefit) GenUpdateSql() string {
//	return fmt.Sprintf("update benefit_tab set name='%s' where id=%d;", p.Name, p.Bid)
//}
//
//func (p Benefit) GenDeleteSql() string {
//	return fmt.Sprintf("update benefit_tab set is_deleted=1 where id=%d;", p.Bid)
//}
//
//func (p Benefit) GenInsertSql() string {
//	return fmt.Sprintf("insert into benefit_tab values(%d,'%s');", p.Bid, p.Name)
//}

func TestUpdate(t *testing.T) {
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

	var p2 = Product{Pid: 123,
		Name: "pro1x",
		Plan: Plan{
			Pid:  345,
			Name: "plan1y",
		},
		Benefits: []Benefit{
			{Bid: 3, Name: "benefit3"},
			{Bid: 2, Name: "benefit2z"},
		},
	}

	fmt.Printf("gen update:\n%s\n", Update(&p1, &p2, "product_tab"))
	fmt.Printf("gen insert:\n%s\n", Update((*Product)(nil), &p2, "product_tab"))
	fmt.Printf("gen delete:\n%s\n", Update(&p1, (*Product)(nil), "product_tab"))

}
