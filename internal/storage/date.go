package storage

import (
	"database/sql/driver"
	"errors"
	"time"
)

/* ---- Date Type ---- */
type Date struct {
	Time  time.Time
	Valid bool
}

func NewDate(t time.Time) Date {
	return Date{Time: t, Valid: true}
}

func (d Date) Add(t time.Duration) Date { return Date{Time: d.Time.Add(t), Valid: d.Valid} }

func (d Date) Format(layout string) string { return d.Time.Format(layout) }

/* ----- Implement JSON Methods ----- */
func (d *Date) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format("2006-01-02") + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	t, err := time.Parse(`"`+"2006-01-02"+`"`, string(data))
	if err != nil {
		return err
	}
	d.Time = t
	d.Valid = true
	return nil
}

/* ----- Implement Date Methods ----- */
func (d Date) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time, nil
}

func (d *Date) Scan(src interface{}) error {
	if src == nil {
		d.Time = time.Time{}
		d.Valid = false
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		d.Time = v
		d.Valid = true
		return nil
	default:
		return errors.New("cannot scan Date")
	}
}

// ---- Utility methods ----
func (d Date) Val() (time.Time, bool) {
	if !d.Valid {
		return time.Time{}, false
	}
	return d.Time, true
}

func (d Date) IsSet() bool {
	return d.Valid
}

func (d Date) Set(t time.Time) {
	d.Time = t
	d.Valid = true
}

func (d Date) Clear() {
	d.Time = time.Time{}
	d.Valid = false
}
