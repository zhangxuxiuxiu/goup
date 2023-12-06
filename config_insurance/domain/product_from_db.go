package domain

import (
	"log"
	"xorm.io/xorm"
)

func BuildProduct(productId uint64, session *xorm.Session) Product {
	var product Product
	exist, err := session.Table("product_tab").Where("is_deleted=0 and id=?", productId).Get(&product)
	if err != nil || !exist {
		log.Fatalf("product with id:%d is not found with err:%v", productId, err)
	}
	product.Plans = buildPlans(productId, session)
	product.PremiumProcessor = buildPremiumProcessor(productId, session)
	product.RiskRules = buildRiskRules(productId, session)
	product.I18ns = buildI18n(productId, session)
	return product
}

func buildI18n(productId uint64, session *xorm.Session) []I18n {
	var i18ns []I18n
	if err := session.Table("i18n_tab").Where("is_deleted=0 and `key` like '%key_?'", productId).Find(&i18ns); err == nil {
		return i18ns
	}
	return nil
}

func buildRiskRules(productId uint64, session *xorm.Session) RiskRule {
	var risk RiskRule
	var condition ConditionRule
	if exist, err := session.Table("uw_condition_tab").Where("is_deleted=0 and product_id=?", productId).Get(&condition); err == nil && exist {
		risk.Condition = &condition
	}
	var scenes []SceneRule
	if err := session.Table("scene_rule_tab").Where("is_deleted=0 and product_id=?", productId).Find(&scenes); err == nil {
		risk.Scenes = scenes
	}
	return risk
}

func buildPremiumProcessor(productId uint64, session *xorm.Session) PremiumProcessor {
	var processor PremiumProcessor
	if exist, err := session.Table("product_premium_config_tab").Where("is_deleted=0 and product_id=?", productId).Get(&processor); err != nil || !exist {
		log.Fatalf("premium processor not found for product id:%d with err:%v", productId, err)
	}
	return processor
}

func buildPlans(productId uint64, session *xorm.Session) []Plan {
	var plans []Plan
	if err := session.Table("plan_tab").Where("is_deleted=0 and product_id=?", productId).Find(&plans); err != nil {
		log.Fatalf("plans not found for product id:%d with err:%v", productId, err)
	}
	for idx := range plans {
		plans[idx].PlanPremium = buildPremium(plans[idx].Id, session)
		plans[idx].Layer = buildLayer(plans[idx].Id, session)
	}
	return plans
}

func buildPremium(planId uint64, session *xorm.Session) PlanPremium {
	var premiums struct {
		ProductId, PlanId, PremiumPriceRuleId, NetPriceRuleId, SumInsuredPriceRuleId uint64
	}
	if exist, err := session.Table("plan_premium_config_tab").Where("is_deleted=0 and plan_id=?", planId).Get(&premiums); err != nil || !exist {
		log.Fatalf("plan premiums not found for plan id:%d with err:%v", planId, err)
	}

	premium := buildPriceRule(premiums.PremiumPriceRuleId, session)
	net := buildPriceRule(premiums.NetPriceRuleId, session)
	sum := buildPriceRule(premiums.SumInsuredPriceRuleId, session)

	return PlanPremium{ProductId: premiums.ProductId, PlanId: premiums.PlanId, Premium: PriceRuleYaml{PriceRule: premium}, Net: PriceRuleYaml{PriceRule: net}, Sum: PriceRuleYaml{PriceRule: sum}}
}

func buildPriceRule(ruleId uint64, session *xorm.Session) PriceRule {
	var generalPriceRule struct {
		Id                     uint64
		RuleType               RuleType
		PriceRuleId            uint64  `yaml:"priceRuleId"`
		RuleDesc               string  `yaml:"rule_desc"`
		LowerBoundType         uint32  `yaml:"LowerBoundType"`
		UpperBoundType         uint32  `yaml:"UpperBoundType"`
		LowerBound             float64 `yaml:"lowerBound" `
		UpperBound             float64 `yaml:"upperBound" `
		TierCount              uint32  `yaml:"tierCount"`
		FixedAmount, FixedRate float64
	}
	if exist, err := session.Table("price_rule_tab").Where("is_deleted=0 and price_rule_id=?", ruleId).Get(&generalPriceRule); err != nil || !exist {
		log.Fatalf("price rules not found for rule id:%d with err:%v", ruleId, err)
	}
	switch generalPriceRule.RuleType {
	case RuleFixedAmount:
		return &PriceRuleFixedAmount{Id: generalPriceRule.Id, LowerBound: generalPriceRule.LowerBound, UpperBound: generalPriceRule.UpperBound,
			LowerBoundType: generalPriceRule.LowerBoundType, UpperBoundType: generalPriceRule.UpperBoundType, TierCount: generalPriceRule.TierCount,
			RuleType: RuleFixedAmount, PriceRuleId: ruleId, RuleDesc: generalPriceRule.RuleDesc, FixedAmount: generalPriceRule.FixedAmount}
	case RuleFixedRate:
		return &PriceRuleFixedRate{Id: generalPriceRule.Id, LowerBound: generalPriceRule.LowerBound, UpperBound: generalPriceRule.UpperBound,
			LowerBoundType: generalPriceRule.LowerBoundType, UpperBoundType: generalPriceRule.UpperBoundType, TierCount: generalPriceRule.TierCount,
			RuleType: RuleFixedRate, PriceRuleId: ruleId, RuleDesc: generalPriceRule.RuleDesc, FixedRate: generalPriceRule.FixedRate}
	case RulePerson:
		var subRules []PriceSubRulePerson
		if err := session.Table("price_rule_person_tab").Where("is_deleted=0 and price_rule_id=?", ruleId).Find(&subRules); err != nil {
			log.Fatalf("plans not found for person rule id:%d with err:%v", ruleId, err)
		}
		return &PriceRulePerson{Id: generalPriceRule.Id, LowerBound: generalPriceRule.LowerBound, UpperBound: generalPriceRule.UpperBound,
			LowerBoundType: generalPriceRule.LowerBoundType, UpperBoundType: generalPriceRule.UpperBoundType, TierCount: generalPriceRule.TierCount,
			RuleType: RulePerson, PriceRuleId: ruleId, RuleDesc: generalPriceRule.RuleDesc, SubRules: subRules}
	case RuleTier:
		var tiers []PriceSubRuleTier
		if err := session.Table("price_rule_tier_tab").Where("is_deleted=0 and price_rule_id=?", ruleId).Find(&tiers); err != nil {
			log.Fatalf("plans not found for person rule id:%d with err:%v", ruleId, err)
		}
		return &PriceRuleTier{Id: generalPriceRule.Id, LowerBound: generalPriceRule.LowerBound, UpperBound: generalPriceRule.UpperBound,
			LowerBoundType: generalPriceRule.LowerBoundType, UpperBoundType: generalPriceRule.UpperBoundType, TierCount: generalPriceRule.TierCount,
			RuleType: RuleTier, PriceRuleId: ruleId, RuleDesc: generalPriceRule.RuleDesc, Tiers: tiers}
	default:
		log.Fatalf("not supported rule type:%d", generalPriceRule.RuleType)
	}
	return nil
}

func buildLayer(planId uint64, session *xorm.Session) *Layer {
	var layer Layer
	layer.PlanId = planId
	if exist, err := session.SQL("select t1.* from layer_tab  t1 join plan_layer_relation_tab t2 on t1.id=t2.layer_id where t2.plan_id=?", planId).
		Get(&layer); err != nil {
		log.Fatalf("layer not found with plan id:%d with err:%v", planId, err)
		return nil
	} else if !exist {
		return nil
	} else {
		layer.Benefits = buildBenefits(layer.Id, session)
		return &layer
	}
}

func buildBenefits(layerId uint64, session *xorm.Session) []Benefit {
	var benefits []Benefit
	if err := session.SQL("select t1.* from benefit_tab t1 join layer_benefit_relation_tab t2 on t1.id=t2.benefit_id where t2.layer_id=?", layerId).
		Find(&benefits); err != nil {
		log.Fatalf("benefits not found with layer id:%d with err:%v", layerId, err)
	}
	return benefits
}
