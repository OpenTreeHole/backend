package utils

import (
	"fmt"
	"testing"

	"github.com/goccy/go-json"
)

type TimeModel struct {
	Time CustomTime `json:"time"`
}

func TestTimeString(t *testing.T) {
	for _, s := range []string{
		"2022-09-09T09:52:55",           // Not RFC3339
		"2022-09-09T09:52:55.123",       // Not RFC3339, millisecond
		"2022-09-09T09:52:55.123456789", // Not RFC3339, nanosecond
		"2022-09-09T09:52:55Z",
		"2022-09-09T09:52:55+08:00",
		"2022-09-09T09:52:55.123Z",
		"2022-09-09T09:52:55.123+08:00",
	} {
		jsonString := fmt.Sprintf(`{"time": "%s"}`, s)
		var model TimeModel
		err := json.Unmarshal([]byte(jsonString), &model)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(model.Time.String())

		jsonBytes, err := json.Marshal(model)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(jsonBytes))
	}
}
