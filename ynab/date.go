package ynab

import (
	"strings"
	"time"
)

// Date is a date-only value that marshals to and from ISO 8601 format (YYYY-MM-DD).
// It embeds time.Time for compatibility with standard time operations.
type Date struct{ time.Time }

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

// NewDate constructs a Date from year, month, and day components.
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}
