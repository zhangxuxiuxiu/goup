package domain

type RiskRule struct {
	Condition *ConditionRule `ipc:"uw_condition_tab"`
	Scenes    []SceneRule    `ipc:"scene_rule_tab"`
}

type ConditionRule struct {
	PartnerId uint64
	ProductId uint64 `ipc:",id"`
	Condition uint32
}

type SceneRule struct {
	PartnerId uint64
	ProductId uint64 `ipc:",id"`
	PlanId    uint64 `ipc:",id"`
	SceneCode string `ipc:",id"`
	RuleCode  string
	Param     string
	RuleIndex uint32 `ipc:",id"`
}
