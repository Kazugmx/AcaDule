package acaduleapi

import (
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02T15:04:05.999999", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ct.Format("2006-01-02T15:04:05.999999") + "\""), nil
}
