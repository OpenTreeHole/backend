package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

type CustomTime struct {
	time.Time
}

var locCST = time.FixedZone("CST", 8*3600)

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	ct.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		ct.Time, err = time.ParseInLocation(`2006-01-02T15:04:05`, s, locCST)
	}
	return err
}

func (ct *CustomTime) UnmarshalText(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	var err error
	ct.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		ct.Time, err = time.ParseInLocation(`2006-01-02T15:04:05`, s, locCST)
	}
	return err
}

func (ct *CustomTime) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal CustomTime value:", value))
	}

	return ct.UnmarshalText(bytes)
}

func (ct CustomTime) MarshalText() ([]byte, error) {
	return json.Marshal(ct.Time)
}

func (ct CustomTime) String() string {
	return ct.Time.Format(time.RFC3339)
}

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.MarshalText()
}
