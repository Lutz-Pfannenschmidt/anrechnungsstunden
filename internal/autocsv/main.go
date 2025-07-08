package autocsv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/andybalholm/crlf"
)

// detectDelimiter detects the delimiter in a given string.
// It assumes the delimiter is the most frequent non-alphanumeric character.
func detectDelimiter(data string) (rune, error) {
	// Count occurrences of each non-alphanumeric character
	charCount := make(map[rune]int)
	for _, char := range data {
		if !strconv.IsPrint(char) || (!unicode.IsLetter(char) && !unicode.IsDigit(char) && !unicode.IsSpace(char)) {
			charCount[char]++
		}
	}

	// Find the character with the highest count
	var delimiter = ','
	maxCount := 0
	if count, ok := charCount[',']; ok {
		maxCount = count
	} else {
		charCount[','] = 0
	}

	for char, count := range charCount {
		if count > maxCount {
			maxCount = count
			delimiter = char
		}
	}

	if maxCount == 0 {
		return ',', fmt.Errorf("could not detect delimiter")
	}

	return delimiter, nil
}

// CSVReader reads a CSV file and auto-detects the delimiter.
// If a specific delimiter is provided, it will use that instead.
// If multiple delimiters are provided, it will use the first one.
func CSVReader(filename string, delim ...rune) (*csv.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	if !scanner.Scan() {
		return nil, fmt.Errorf("file is empty")
	}
	firstLine := scanner.Text()

	var delimiter rune
	if len(delim) > 0 {
		delimiter = delim[0]
	} else {
		delimiter, err = detectDelimiter(firstLine)
		if err != nil {
			return nil, err
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("could not seek to beginning of file: %w", err)
	}

	csvReader := csv.NewReader(crlf.NewReader(file))
	csvReader.LazyQuotes = true
	csvReader.Comma = delimiter

	return csvReader, nil
}
