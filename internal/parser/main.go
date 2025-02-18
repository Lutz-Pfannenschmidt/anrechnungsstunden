package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Lutz-Pfannenschmidt/stunden-berechner/internal/csv"
)

func ParseFile(path string) (map[string][2]float64, error) {
	lines, err := csv.ReadAnyFileToCSV(path)
	if err != nil {
		return nil, err
	}

	var result = map[string][2]float64{}

	var currName string
	var currYear [2][]float64
	var foundTable bool

	for _, line := range *lines {
		line = strings.TrimSpace(line)
		if isNameLine(line) && !foundTable {
			currName = strings.ReplaceAll(line, ",", "")
			currYear = [2][]float64{}
			foundTable = false
		} else if isHeaderLine(line) {
			foundTable = true
		} else if foundTable && currName != "" {
			if strings.Contains(line, "Ferien") {
				continue
			}

			values := strings.Split(line, ",")

			if values[0] == "" {
				avg1 := avg(currYear[0])
				avg2 := avg(currYear[1])
				result[currName] = [2]float64{avg1, avg2}
				foundTable = false
				continue
			}

			c, err := strconv.ParseFloat(values[4], 64)
			if err != nil {
				return nil, err
			}

			weeks := 1
			if strings.Contains(values[0], "-") {
				var from, to int
				fmt.Sscanf(values[0], "%d-%d", &from, &to)
				weeks = to - from + 1
			}

			var period uint
			fmt.Sscanf(values[2], "%d", &period)

			sem := 0
			if period >= 4 {
				sem = 1
			}

			for i := 0; i < weeks; i++ {
				currYear[sem] = append(currYear[sem], c)
			}
		}

	}

	return result, nil
}

func isHeaderLine(line string) bool {
	return strings.Contains(line, "Woche") && strings.Contains(line, "Periode") && strings.Contains(line, "Soll")
}

// // expandLines expands lines if they contain a range of weeks.
// func expandLines(lines []string) []string {
// 	var expanded []string
// 	for _, line := range lines {
// 		values := strings.Split(line, ",")
// 		if strings.Contains(values[0], "-") {
// 			var from, to int
// 			fmt.Sscanf(values[0], "%d-%d", &from, &to)
// 			var fromDate, toDate string
// 			fmt.Sscanf(values[1], "%s-%s", &fromDate, &toDate)
// 			fromD := date.MustParseDate(fromDate)
// 			toD := date.MustParseDate(toDate)
// 			for i := from; i <= to; i++ {
// 				expanded = append(expanded, fmt.Sprintf("%d,%d.%d.-%d.%d.,%s,", i, strings.Join(values[2:], ",")))
// 			}
// 		} else {
// 			expanded = append(expanded, line)
// 		}
// 	}
// 	return expanded
// }

// isNameLine checks if a line is a name line by checking if it contains only "a-zA-Z_-," and has at least 3 letters.
func isNameLine(line string) bool {
	line = strings.ReplaceAll(line, ",", "")
	line = strings.ReplaceAll(line, "_", "")
	line = strings.ReplaceAll(line, "-", "")
	if len(line) < 3 {
		return false
	}

	for i := 0; i < len(line); i++ {
		if !isLetter(line[i]) {
			return false
		}
	}

	return true
}

func isLetter(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char == 'ä' || char == 'ö' || char == 'ü' || char == 'ß'
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
