package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	date int
}

// ParseDate creates a new Date object from a string where the date is in the format "d.m.".
func ParseDate(date string) (*Date, error) {
	parsed := strings.Split(date, ".")
	if len(parsed) < 2 {
		return nil, fmt.Errorf("invalid date format: %s", date)
	}
	d, err := strconv.Atoi(parsed[0])
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s", date)
	}
	m, err := strconv.Atoi(parsed[1])
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %s", date)
	}

	return &Date{date: d + m*100}, nil

}

func MustParseDate(date string) *Date {
	d, err := ParseDate(date)
	if err != nil {
		panic(err)
	}
	return d
}

// FromInt creates a new Date object from an integer in the format "ddmm".
func FromInt(date int) *Date {
	return &Date{date: date}
}

// GetInt returns the date as an integer in the format "ddmm".
func (d *Date) GetInt() int {
	return d.date
}

func (d *Date) String() string {
	return fmt.Sprintf("%02d.%02d.", d.date%100, d.date/100)
}

// Compare compares two Date objects and returns -1 if d is before other, 1 if d is after other, and 0 if they are equal.
func (d *Date) Compare(other *Date) int {
	if d.date < other.date {
		return -1
	}
	if d.date > other.date {
		return 1
	}
	return 0
}

// DaysUntil returns the number of days between d and other. If other is before d, other is treated as beeing in the next year.
func (d *Date) DaysUntil(other *Date) int {
	tmp := time.Date(2024, time.Month(d.date/100), d.date%100, 0, 0, 0, 0, time.UTC)
	tmp2 := time.Date(2024, time.Month(other.date/100), other.date%100, 0, 0, 0, 0, time.UTC)

	if tmp2.Before(tmp) {
		tmp2 = tmp2.AddDate(1, 0, 0)
	}
	return int(tmp2.Sub(tmp).Hours() / 24)
}

// GetMonth returns the month of the date.
func (d *Date) GetMonth() int {
	return d.date / 100
}

// GetDay returns the day of the date.
func (d *Date) GetDay() int {
	return d.date % 100
}
