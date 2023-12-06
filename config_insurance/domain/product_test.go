package domain_test

import (
	"fmt"
	"github.com/go_practice/config_insurance/dbops"
	"github.com/go_practice/config_insurance/domain"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"testing"
)

func TestDomainYaml(t *testing.T) {
	bytes, err := os.ReadFile("./domain.yaml")
	if err != nil {
		t.Fatalf("file not found")
	}
	var product domain.Product
	if err := yaml.Unmarshal(bytes, &product); err != nil {
		t.Fatalf("Product Unmarshal failed with:%s", err.Error())
	}

	b1, _ := yaml.Marshal(product.Plans[0].PlanPremium.Premium)
	fmt.Printf("product premiums:%s\n", b1)
	//fmt.Printf("gen sql:%s\n", product.GenInsertSql())

	//b, _ := yaml.Marshal(struct{ domain.PriceRule }{PriceRule: &domain.PriceRuleFixedAmount{Type: domain.RuleFixedAmount, RuleId: 1, FixedAmount: 2}})
	//var priceRule domain.PriceRule = &domain.PriceRuleFixedAmount{Type: domain.RuleFixedAmount, RuleId: 1, FixedAmount: 2}
	//b, _ := yaml.Marshal(priceRule)
	//fmt.Printf("yaml plan premium:%s\n", b)
}

func TestInterfaceImplements(t *testing.T) {
	var layer = domain.Layer{}
	typPlan := reflect.TypeOf(layer)
	genDelType := reflect.TypeOf((*dbops.GenDelete)(nil)).Elem()
	t.Logf("does Layer implements GenDelete:%v\n", typPlan.Implements(genDelType))
	t.Logf("does *Layer implements GenDelete:%v\n", reflect.PtrTo(typPlan).Implements(genDelType))

	valPlan := reflect.ValueOf(layer)
	_, ok := valPlan.Interface().(dbops.GenDelete)
	t.Logf("does Layer cast to GenDelete:%v\n", ok)
	if valPlan.CanAddr() {
		_, ok := valPlan.Addr().Interface().(dbops.GenDelete)
		t.Logf("does *Layer cast to GenDelete:%v\n", ok)
	} else {
		t.Logf("Layer var can't addr\n")
	}

	pvalPlan := reflect.ValueOf(&layer)
	t.Logf("*Layer.Elem can addr:%v\n", pvalPlan.Elem().CanAddr())
	_, ok = pvalPlan.Interface().(dbops.GenDelete)
	t.Logf("does *Layer cast to GenDelete:%v\n", ok)
	t.Logf("pvalPlan.Elem().Kind()=%v", pvalPlan.Elem().Kind())
	//pvalPlan.Elem().
	//t.Logf("does *Layer.ProductId can addr:%v\n", pvalPlan.Elem().Field(0).CanAddr())

}

func TestEnum(*testing.T) {
	fmt.Printf("%v", domain.RuleFixedAmount)
}
