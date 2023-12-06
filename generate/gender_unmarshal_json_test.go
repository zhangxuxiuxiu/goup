package gen

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	strs := []string{"Unknown", "Male", "Female"}
	//gender := new(Gender)
	var tmp = struct {
		Gen Gender `json:"gender"`
	}{}
	for _, str := range strs {
		if err := json.Unmarshal(([]byte)(fmt.Sprintf("{\"gender\":\"%s\"}", str)), &tmp); err != nil {
			t.Fatalf("fail with :%s", err.Error())
		}
		assert.Equal(t, tmp.Gen, _Gender_name_to_values[str])
	}
}
