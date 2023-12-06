package domain

import (
	"fmt"
	"github.com/go_practice/config_insurance/dbops"
	"strings"
)

// id, code, name, full_name,`desc`, main_type, sub_type,period_category,
// period_type,period_type_value,broker_id,broker_product_id,
// renew,status,default_plan_id,end_time,insurer_id,insurer_name
type Product struct {
	Id               uint64 `yaml:"productId" xorm:"id" ipc:",id"`
	Code             string
	Name             string
	FullName         string `yaml:"fullName"`
	Desc             string
	MainType         uint32 `yaml:"mainType"`
	SubType          uint32 `yaml:"subType"`
	PeriodCategory   uint32 `yaml:"periodCategory"`
	PeriodType       uint32 `yaml:"periodType"`
	PeriodTypeValue  string `yaml:"periodTypeValue"`
	BrokerId         uint32 `yaml:"brokerId"`
	BrokerProductId  string `yaml:"brokerProductId"`
	Renew            uint32
	Status           uint32
	DefaultPlanId    uint64           `yaml:"defaultPlanId"`
	EndTime          string           `yaml:"endTime"`
	InsurerId        uint64           `yaml:"insurerId"`
	InsurerName      string           `yaml:"insurerName"`
	PremiumProcessor PremiumProcessor `yaml:"premiumProcessor" xorm:"-" ipc:"product_premium_config_tab"`
	Plans            []Plan           `xorm:"-" ipc:"plan_tab"`
	RiskRules        RiskRule         `xorm:"-" yaml:"riskRules"`
	I18ns            []I18n           `xorm:"-" ipc:"i18n_tab"`
}

type I18n struct {
	Key      string `ipc:",id"`
	Language string `ipc:",id"`
	Wording  string
}

type PremiumProcessor struct {
	ProductId        uint64 `yaml:"productId" ipc:",id"`
	PremiumProcessor string `yaml:"premiumProcessor"`
}

// id, product_id,`code`,name,`desc`,broker_id,
// broker_plan_id,period,max_multiple_insured
type Plan struct {
	ProductId          uint64 `yaml:"productId" ipc:",id"`
	Id                 uint64 `ipc:",id"`
	Code               string
	Name               string
	Desc               string
	BrokerId           uint64 `yaml:"brokerId"`
	BrokerPlanId       string `yaml:"brokerPlanId"`
	Period             string
	MaxMultipleInsured uint64
	PlanPremium        PlanPremium `yaml:"planPremium" xorm:"-" ipc:"plan_premium_config_tab"`
	Layer              *Layer      `xorm:"-" ipc:"layer_tab"`
}

type Layer struct {
	PlanId   uint64 `ipc:"-" yaml:"planId"`
	Id       uint64 `ipc:",id"`
	Code     string
	Benefits []Benefit `xorm:"-"`
}

func (p *Layer) GenDeleteSql() string {
	if p == nil {
		return ""
	}
	return strings.Join([]string{
		fmt.Sprintf("update plan_layer_relation_tab set is_deleted=1 where id=%d and plan_id=%d", p.Id, p.PlanId),
		fmt.Sprintf("update layer_tab set is_deleted=1 where id=%d", p.Id),
		fmt.Sprintf("update layer_benefit_relation_tab set is_deleted=1 where layer_id=%d", p.Id),
		fmt.Sprintf("update benefit_tab set is_deleted=1 where id in (select  a.benefit_id from layer_benefit_relation_tab a where layer_id=%d);", p.Id),
	}, "\n")
}

/*
mysql> select * from plan_layer_relation_tab ;
+----+---------------------+---------------------+---------+-------------+---------------------+---------------------+---------+------------+--------------+-----------+
| id | plan_id             | layer_id            | is_main | is_optional | gmt_created         | gmt_modified        | version | is_deleted | exclusion_id | layer_seq |
+----+---------------------+---------------------+---------+-------------+---------------------+---------------------+---------+------------+--------------+-----------+
|  1 | 1569261920811221252 | 1529743333629368320 |       1 |           2 | 2021-12-30 20:39:28 | 2021-12-30 20:39:28 |       0 |          0 |            0 |         0 |

mysql> select * from layer_tab;
+---------------------+--------+---------+------+-----------+---------------------+---------------------+---------+------------+-------+
| id                  | code   | name    | desc | main_type | gmt_created         | gmt_modified        | version | is_deleted | extra |
+---------------------+--------+---------+------+-----------+---------------------+---------------------+---------+------------+-------+
| 1529743333629368320 | l_01   | default |      |         1 | 2021-12-30 20:39:29 | 2021-12-30 20:39:29 |       0 |          0 |       |
+---------------------+--------+---------+------+-----------+---------------------+---------------------+---------+------------+-------+
*/
func (p *Layer) GenInsertSql() string {
	if p == nil {
		return ""
	}
	layerSql := fmt.Sprintf("insert into layer_tab (id,code,name,main_type) values (%d,\"%s\",\"default\",1); ", p.Id, p.Code)
	layerRelateSql := fmt.Sprintf("insert into plan_layer_relation_tab(plan_id, layer_id, is_main, is_optional) values (%d,%d,1,2); ", p.PlanId, p.Id)

	sqls := make([]string, len(p.Benefits)*2+2)
	sqls[0] = layerSql
	sqls[1] = layerRelateSql
	for idx, b := range p.Benefits {
		sqls[idx*2+2] = dbops.Insert(&b, "benefit_tab") //.GenInsertSql()
		sqls[idx*2+3] = fmt.Sprintf("insert into layer_benefit_relation_tab(layer_id, benefit_id, is_optional) values (%d,%d,2);", p.Id, b.Id)
	}

	return strings.Join(sqls, "\n")
}

type Benefit struct {
	Id             uint64 `ipc:",id"`
	Code           string
	Name           string
	SumInsuredType uint64 `yaml:"sumInsuredType"`
	SumInsuredDesc string `yaml:"sumInsuredDesc"`
}

type PlanPremium struct {
	ProductId         uint64 `yaml:"productId" ipc:",id"`
	PlanId            uint64 `yaml:"planId" ipc:",id"`
	Premium, Net, Sum PriceRuleYaml
}

/*
mysql> select * from plan_tab where product_id=1385277174484245481;
+---------------------+---------------------+--------+------+---------------------+-----------+---------------------+---------------------+---------------------+---------+------------+--------+
| id                  | product_id          | code   | name | desc                | broker_id | broker_plan_id      | gmt_created         | gmt_modified        | version | is_deleted | period |
+---------------------+---------------------+--------+------+---------------------+-----------+---------------------+---------------------+---------------------+---------+------------+--------+
| 1569261920811221252 | 1385277174484245481 | PA_P01 |      | PLAN-PA-A,PLAN-PA-D |         2 | PLAN-PA-A,PLAN-PA-D | 2021-12-30 20:38:55 | 2022-08-29 02:42:54 |       0 |          0 | 1M,12M |
+---------------------+---------------------+--------+------+---------------------+-----------+---------------------+---------------------+---------------------+---------+------------+--------+

mysql> select * from product_tab where id=1385277174484245481;
+---------------------+------+------------------------------+------------------------------+------+----------+-----------+----------+-----------------+-------------+-------------------+-----------------+-----------+-------------------+--------------+--------------+----------------+-------+--------+---------------------+---------------------+---------------------+---------------------+---------+------------+------------+
| id                  | code | name                         | full_name                    | desc | logo_url | main_type | sub_type | period_category | period_type | period_type_value | default_plan_id | broker_id | broker_product_id | insurer_code | insurer_name | require_active | renew | status | start_time          | end_time            | gmt_created         | gmt_modified        | version | is_deleted | insurer_id |
+---------------------+------+------------------------------+------------------------------+------+----------+-----------+----------+-----------------+-------------+-------------------+-----------------+-----------+-------------------+--------------+--------------+----------------+-------+--------+---------------------+---------------------+---------------------+---------------------+---------+------------+------------+
| 1385277174484245481 | P08  | Personal Accident Protection | Personal Accident Protection |      |          |         3 |      103 |               1 |           3 | 1,12              |               0 |         2 | PROD-PA-001       |              |              |              0 |     2 |      1 | 2021-12-09 07:03:25 | 2099-11-11 03:11:20 | 2021-12-30 20:38:55 | 2022-08-29 02:32:50 |       0 |          0 |          0 |
+---------------------+------+------------------------------+------------------------------+------+----------+-----------+----------+-----------------+-------------+-------------------+-----------------+-----------+-------------------+--------------+--------------+----------------+-------+--------+---------------------+---------------------+---------------------+---------------------+---------+------------+------------+
*/
func (p *PlanPremium) GenInsertSql() string {
	premiumConf := fmt.Sprintf("insert into plan_premium_config_tab (product_id, plan_id, premium_price_rule_id, net_price_rule_id, sum_insured_price_rule_id) "+
		" values (%d,%d,%d,%d,%d);", p.ProductId, p.PlanId, p.Premium.RId(), p.Net.RId(), p.Sum.RId())
	return fmt.Sprintf("%s\n%s\n%s\n%s", premiumConf, dbops.Insert(&p.Premium, "price_rule_tab"), dbops.Insert(&p.Net, "price_rule_tab"), dbops.Insert(&p.Sum, "price_rule_tab"))
}
