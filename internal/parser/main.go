package parser

import (
	"anrechnungsstundenberechner/internal/autocsv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DATE_LAYOUT = "2.1.2006" // DD.MM.YYYY

type ParseResult struct {
	Result         map[string][2]float64 `json:"result"`
	NameCollisions map[string][]string   `json:"name_collisions"`
}

func isMonday(split time.Time) bool {
	return split.Weekday() == time.Monday
}

func parseAndValidateSplitDate(splitDate string) (*time.Time, error) {
	err_invalid := fmt.Errorf("split_date must be in MM.DD.YYYY format: %s", splitDate)

	if len(splitDate) != 10 {
		return nil, err_invalid
	}

	if splitDate[2] != '.' || splitDate[5] != '.' {
		return nil, err_invalid
	}

	t, err := time.Parse(DATE_LAYOUT, splitDate)
	if err != nil {
		return nil, err_invalid
	}

	d := t.Weekday()

	if d != time.Monday {
		return nil, fmt.Errorf("split_date must be a monday: %s", splitDate)
	}

	return &t, nil
}

func ParseFile(path string, fromYear int, SplitDateStr string) (res *ParseResult, err error) {
	splitDate, err := parseAndValidateSplitDate(SplitDateStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing split date: %w", err)
	}

	data, err := autocsv.ReadAnyFileToCSV(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if !isMonday(*splitDate) {
		return nil, fmt.Errorf("invalid split date: %s, must be a Monday", splitDate.Format(DATE_LAYOUT))
	}

	fromDate, _, err := findSchoolYear(data, fromYear)
	if err != nil {
		return nil, fmt.Errorf("error finding school year: %w", err)
	}

	res = &ParseResult{
		Result:         map[string][2]float64{},
		NameCollisions: map[string][]string{},
	}

	parts := strings.SplitSeq(data, "Stundenplan / Werte")

	for part := range parts {
		if part == "" {
			continue
		}

		name := uniqueName(strings.TrimSpace(strings.ReplaceAll(strings.Split(part, "Jahresmittelwert")[0], ",", "")), res)

		foundTable := false
		yearData := [2][]float64{}

		lines := strings.SplitSeq(part, "\n")
		for l := range lines {
			l = strings.TrimSpace(l)

			if strings.HasPrefix(l, "Woche,") {
				foundTable = true
				continue
			}

			if !foundTable || strings.Contains(l, "Ferien") {
				continue
			}

			values := strings.Split(l, ",")

			if values[0] == "" {
				avg1 := avg(yearData[0])
				avg2 := avg(yearData[1])
				res.Result[name] = [2]float64{avg1, avg2}

				foundTable = false
				continue
			}

			c, err := strconv.ParseFloat(values[4], 64)
			if err != nil {
				fmt.Println("Error parsing value:", values[4])
				return nil, err
			}

			weeks := 1
			if strings.Contains(values[0], "-") {
				var from, to int
				_, err = fmt.Sscanf(values[0], "%d-%d", &from, &to)
				if err != nil {
					fmt.Println("Error parsing weeks:", values[0])
					return nil, err
				}
				weeks = to - from + 1
			}

			var fromMonth, toMonth time.Month
			var fromDay, toDay int
			n, err := fmt.Sscanf(values[1], "%d.%d.-%d.%d.", &fromDay, &fromMonth, &toDay, &toMonth)
			if err != nil {
				fmt.Println("Error parsing date:", values[1])
				return nil, err
			} else if n != 4 {
				fmt.Println("Error parsing date:", values[1])
				return nil, fmt.Errorf("invalid date format: %s", values[1])
			}

			sem := 0
			if isDayMonthAfterSplitDate(toDay, toMonth, splitDate, fromDate) {
				sem = 1
			}

			for range weeks {
				yearData[sem] = append(yearData[sem], c)
			}

		}

	}

	for key := range res.Result {
		if _, exists := res.Result[key+"_NAME_COLLISION_1"]; exists {
			res.Result[key+"_NAME_COLLISION_0"] = res.Result[key]
			delete(res.Result, key)
		}
	}

	return
}

// isDayMonthFirstSemester checks if the given day and month are in the first semester
// splitDate is the date when the semesters are split (everything before that date is first semester).
// yearStart is the date when the school year starts. Anything before that must be in the second YEAR but not necessarily in the second semester.
func isDayMonthAfterSplitDate(day int, month time.Month, splitDate, yearStart *time.Time) bool {
	year := yearStart.Year()
	if month < yearStart.Month() || (month == yearStart.Month() && day < yearStart.Day()) {
		year++
	}

	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return t.After(*splitDate)
}

func uniqueName(name string, res *ParseResult) string {
	newName := name
	i := 1
	for {
		if _, exists := res.Result[newName]; exists {
			if len(res.NameCollisions[name]) == 0 {
				res.NameCollisions[name] = []string{name + "_NAME_COLLISION_0"}
			}
			newName = fmt.Sprintf("%s_NAME_COLLISION_%d", name, i)
			if _, exists := res.Result[newName]; !exists {
				res.NameCollisions[name] = append(res.NameCollisions[name], newName)
				break
			}
			i++
		} else {
			break
		}
	}
	return newName
}

func findSchoolYear(data string, startYear int) (*time.Time, *time.Time, error) {
	r := strings.NewReplacer(" ", "", "Wochenwerte", "")
	from := time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)
	to := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Wochenwerte") {
			l := r.Replace(line)
			parts := strings.Split(l, "-")
			if len(parts) != 2 {
				continue
			}

			fromDate := strings.TrimSpace(parts[0]) + strconv.Itoa(startYear)
			toDate := strings.TrimSpace(parts[1]) + strconv.Itoa(startYear+1)

			f, err := time.Parse(DATE_LAYOUT, fromDate)
			if err != nil {
				continue
			}

			t, err := time.Parse(DATE_LAYOUT, toDate)
			if err != nil {
				continue
			}

			if f.Before(from) {
				from = f
			}
			if t.After(to) {
				to = t
			}

		}
	}

	if from.Equal(time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)) {
		return nil, nil, fmt.Errorf("no valid date found in lines")
	}
	if to.IsZero() {
		return nil, nil, fmt.Errorf("no valid date found in lines")
	}
	if from.After(to) {
		return nil, nil, fmt.Errorf("from date is after to date")
	}

	return &from, &to, nil
}

// avg calculates the average of an array of uints.
// If the array is empty, 0 is returned.
func avg(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	sum := 0.0
	l := float64(len(arr))
	for _, f := range arr {
		sum += f
	}
	return sum / l
}
