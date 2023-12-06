package domain

import (
	"fmt"
	"reflect"
)

//go:generate go-enum -trimprefix=Rule  -all=false -yaml -string -type=RuleType
type RuleType uint8

const (
	RuleTier        RuleType = 1
	RuleFixedAmount RuleType = 2
	RuleFixedRate   RuleType = 3
	RuleTenure      RuleType = 4
	RulePerson      RuleType = 5
)

type PriceRule interface {
	RId() uint64
	RType() RuleType
	fmt.Stringer
}

var priceRules = map[RuleType]reflect.Type{}

func init() {
	for _, rule := range []PriceRule{&PriceRulePerson{}, &PriceRuleFixedAmount{}, &PriceRuleFixedRate{}, &PriceRuleTier{}} {
		priceRules[rule.RType()] = reflect.TypeOf(rule).Elem()
	}
}

// id, price_rule_id, describ, min_age, max_age, period, period_type,fixed_amount,age_type
type PriceSubRulePerson struct {
	Id          uint64 `ipc:"-"`
	PriceRuleId uint64 `yaml:"ruleId" ipc:",id"`
	Describ     string
	MinAge      uint16  `yaml:"minAge" ipc:",id"`
	MaxAge      uint16  `yaml:"maxAge" ipc:",id"`
	Period      uint16  `ipc:",id"`
	PeriodType  uint16  `yaml:"periodType" ipc:",id"`
	AgeType     uint64  `yaml:"ageType" ipc:",id"`
	FixedAmount float64 `yaml:"fixedAmount"`
}

// id,price_rule_id,rule_type,rule_desc,lower_bound_type,upper_bound_type,tier_count
type PriceRulePerson struct {
	Id             uint64               `ipc:"-"`
	RuleType       RuleType             `yaml:"ruleType"`
	PriceRuleId    uint64               `yaml:"priceRuleId" ipc:",id"`
	RuleDesc       string               `yaml:"ruleDesc"`
	LowerBoundType uint32               `yaml:"lowerBoundType"`
	UpperBoundType uint32               `yaml:"upperBoundType"`
	LowerBound     float64              `yaml:"lowerBound" ipc:",id"`
	UpperBound     float64              `yaml:"upperBound" ipc:",id"`
	TierCount      uint32               `yaml:"tierCount"`
	SubRules       []PriceSubRulePerson `yaml:"subRules" ipc:"price_rule_person_tab"`
}

func (p *PriceRulePerson) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *PriceRulePerson) RType() RuleType {
	return RulePerson
}

func (p *PriceRulePerson) RId() uint64 {
	return p.PriceRuleId
}

// price_rule_id,tier_seq, lower_bound, upper_bound,lower_bound_type, upper_bound_type,amount
type PriceSubRuleTier struct {
	Id             uint64  `ipc:"-"`
	PriceRuleId    uint64  `yaml:"ruleId" ipc:",id"`
	TierSeq        uint16  `yaml:"tierSeq" ipc:",id"`
	LowerBound     float64 `yaml:"lowerBound" ipc:",id"`
	UpperBound     float64 `yaml:"upperBound" ipc:",id"`
	LowerBoundType uint16  `yaml:"lowerBoundType" ipc:",id"`
	UpperBoundType uint16  `yaml:"upperBoundType" ipc:",id"`
	Amount         float64
}

type PriceRuleTier struct {
	Id             uint64             `ipc:"-"`
	RuleType       RuleType           `yaml:"ruleType"`
	PriceRuleId    uint64             `yaml:"priceRuleId" ipc:",id"`
	RuleDesc       string             `yaml:"ruleDesc"`
	LowerBoundType uint32             `yaml:"lowerBoundType"`
	UpperBoundType uint32             `yaml:"upperBoundType"`
	LowerBound     float64            `yaml:"lowerBound" ipc:",id"`
	UpperBound     float64            `yaml:"upperBound" ipc:",id"`
	TierCount      uint32             `yaml:"tierCount"`
	Tiers          []PriceSubRuleTier `yaml:"tiers" ipc:"price_rule_tier_tab"`
}

func (p *PriceRuleTier) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *PriceRuleTier) RType() RuleType {
	return RuleTier
}

func (p *PriceRuleTier) RId() uint64 {
	return p.PriceRuleId
}

type PriceRuleFixedAmount struct {
	Id             uint64   `ipc:"-"`
	RuleType       RuleType `yaml:"ruleType"`
	PriceRuleId    uint64   `yaml:"priceRuleId" ipc:",id"`
	RuleDesc       string   `yaml:"ruleDesc"`
	LowerBoundType uint32   `yaml:"lowerBoundType"`
	UpperBoundType uint32   `yaml:"upperBoundType"`
	LowerBound     float64  `yaml:"lowerBound" ipc:",id"`
	UpperBound     float64  `yaml:"upperBound" ipc:",id"`
	TierCount      uint32   `yaml:"tierCount"`
	FixedAmount    float64  `yaml:"fixedAmount"`
}

func (p *PriceRuleFixedAmount) String() string {
	return fmt.Sprintf("%v", *p)
}

func (p *PriceRuleFixedAmount) RType() RuleType {
	return RuleFixedAmount
}

func (p *PriceRuleFixedAmount) RId() uint64 {
	return p.PriceRuleId
}

func (p *PriceRuleFixedRate) String() string {
	return fmt.Sprintf("%v", *p)
}

type PriceRuleFixedRate struct {
	Id             uint64   `ipc:"-"`
	RuleType       RuleType `yaml:"ruleType"`
	PriceRuleId    uint64   `yaml:"priceRuleId" ipc:",id"`
	RuleDesc       string   `yaml:"ruleDesc"`
	LowerBoundType uint32   `yaml:"lowerBoundType"`
	UpperBoundType uint32   `yaml:"upperBoundType"`
	LowerBound     float64  `yaml:"lowerBound" ipc:",id"`
	UpperBound     float64  `yaml:"upperBound" ipc:",id"`
	TierCount      uint32   `yaml:"tierCount"`
	FixedRate      float64  `yaml:"fixedRate"`
}

func (p *PriceRuleFixedRate) RType() RuleType {
	return RuleFixedRate
}

func (p *PriceRuleFixedRate) RId() uint64 {
	return p.PriceRuleId
}

var (
	_ PriceRule = &PriceRulePerson{}
	_ PriceRule = &PriceRuleFixedAmount{}
	_ PriceRule = &PriceRuleFixedRate{}
)

// PriceRuleYaml  a PriceRule wrapper to be parsed from yaml file
type PriceRuleYaml struct {
	PriceRule `ipc:"price_rule_tab"`
}

func (p *PriceRuleYaml) String() string {
	return fmt.Sprintf("%v", p.PriceRule.String())
}

func (p *PriceRuleYaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	ruleType := struct {
		Type RuleType `yaml:"ruleType"`
	}{}
	if err := unmarshal(&ruleType); err != nil {
		return err
	}
	reflect.ValueOf(&p.PriceRule).Elem().Set(reflect.New(priceRules[ruleType.Type]))

	return unmarshal(p.PriceRule)
}

func (p PriceRuleYaml) MarshalYAML() (interface{}, error) {
	//if marshal PriceRule as string here, marshaled text will be regarded as multiline string,
	//which can't be unmarshalled back to Product
	return p.PriceRule, nil
}
