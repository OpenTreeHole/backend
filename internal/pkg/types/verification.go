package types

import (
	"fmt"
	"strconv"
	"strings"
)

type VerificationCode string

func (v *VerificationCode) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "null" {
		return nil
	}

	number, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	*v = VerificationCode(fmt.Sprintf("%06d", number))
	return nil
}

func (v *VerificationCode) UnmarshalText(data []byte) error {
	s := strings.Trim(string(data), `"`)
	// Ignore null, like in the main JSON package.
	if s == "" {
		return nil
	}

	*v = VerificationCode(s)
	return nil
}
