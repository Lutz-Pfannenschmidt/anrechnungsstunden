package parser

import (
	"testing"
	"time"
)

func TestIsDayMonthSecondSemester(t *testing.T) {
	splitDate, err := time.Parse(DATE_LAYOUT, "01.09.2023")
	if err != nil {
		t.Fatalf("error parsing split date: %v", err)
	}
	splitDate2, err := time.Parse(DATE_LAYOUT, "12.01.2024")
	if err != nil {
		t.Fatalf("error parsing split date: %v", err)
	}
	yearStart, err := time.Parse(DATE_LAYOUT, "01.08.2023")
	if err != nil {
		t.Fatalf("error parsing year start date: %v", err)
	}

	tests := []struct {
		day       int
		month     time.Month
		want      bool
		splitDate *time.Time
	}{
		{1, time.January, true, &splitDate},
		{1, time.February, true, &splitDate},
		{1, time.March, true, &splitDate},
		{1, time.April, true, &splitDate},
		{1, time.May, true, &splitDate},
		{1, time.June, true, &splitDate},
		{1, time.July, true, &splitDate},
		{1, time.August, false, &splitDate},
		{2, time.August, false, &splitDate},
		{31, time.August, false, &splitDate},
		{1, time.September, true, &splitDate},
		{2, time.September, true, &splitDate},
		{1, time.October, true, &splitDate},
		{1, time.November, true, &splitDate},
		{1, time.December, true, &splitDate},
		{1, time.December, false, &splitDate2},
		{31, time.December, false, &splitDate2},
		{1, time.January, false, &splitDate2},
		{12, time.January, true, &splitDate2},
		{1, time.February, true, &splitDate2},
	}

	for _, test := range tests {
		got := isDayMonthAfterSplitDate(test.day, test.month, test.splitDate, &yearStart)
		if got != test.want {
			t.Errorf("isDayMonthSecondSemester(%d, %v, %v) = %v; want %v", test.day, test.month, test.splitDate.Format(DATE_LAYOUT), got, test.want)
		}
	}
}
