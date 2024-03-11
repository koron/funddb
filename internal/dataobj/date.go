package dataobj

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

func NewDate(year int, month time.Month, day int) Date {
	ti := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return DateFromTime(ti)
}

func DateFromTime(ti time.Time) Date {
	return Date{ Year: ti.Year(), Month: int(ti.Month()), Day: ti.Day()}
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

var _ driver.Valuer = Date{}

func (d Date) Value() (driver.Value, error) {
	if d.Year == 0 && d.Month == 0 && d.Day == 0 {
		return nil, nil
	}
	return d.String(), nil
}

var _ sql.Scanner = (*Date)(nil)

func (d *Date) Scan(src any) error {
	ti, err := time.Parse(time.DateOnly, fmt.Sprint(src))
	if err != nil {
		return err
	}
	*d = DateFromTime(ti)
	return nil
}
