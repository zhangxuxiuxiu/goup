package util

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PremiumStrategy uint8

const (
	InitiatePremium PremiumStrategy = 1
	TrialPremium    PremiumStrategy = 2
	ConfirmPremium  PremiumStrategy = 3
)

var premiumStrategies = map[string]PremiumStrategy{
	"initial": InitiatePremium,
	"trial":   TrialPremium,
	"confirm": ConfirmPremium,
}

func (a *PremiumStrategy) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("Strategy should be a string, got %s", b)
	}
	strategy, ok := premiumStrategies[s]
	if ok {
		*a = strategy
		return nil
	}
	return fmt.Errorf("invalid premium strategy:%s", string(b))
}

func TestUnmarshalJSON(t *testing.T) {
	strs := []string{"initial", "trial", "confirm"}
	var tmp = struct {
		S PremiumStrategy `json:"strategy"`
	}{}
	for _, str := range strs {

		if err := json.Unmarshal(([]byte)(fmt.Sprintf("{\"strategy\": \"%s\"}", str)), &tmp); err != nil {
			t.Fatalf("fail with :%s", err.Error())
		}
		assert.Equal(t, tmp.S, premiumStrategies[str])
	}
}
