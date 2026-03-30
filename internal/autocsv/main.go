package autocsv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"

	"github.com/andybalholm/crlf"
)

// UTF-8 BOM bytes
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

// bomSkippingReader wraps an io.Reader and skips the UTF-8 BOM if present
type bomSkippingReader struct {
	reader     io.Reader
	bomChecked bool
}

func (r *bomSkippingReader) Read(p []byte) (int, error) {
	if !r.bomChecked {
		r.bomChecked = true
		// Read first 3 bytes to check for BOM
		bom := make([]byte, 3)
		n, err := io.ReadFull(r.reader, bom)
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			return 0, err
		}
		// Check if it's a UTF-8 BOM
		if n == 3 && bom[0] == utf8BOM[0] && bom[1] == utf8BOM[1] && bom[2] == utf8BOM[2] {
			// BOM found and skipped, continue reading normally
			return r.reader.Read(p)
		}
		// No BOM, we need to return the bytes we already read
		copy(p, bom[:n])
		if n < len(p) {
			m, err := r.reader.Read(p[n:])
			return n + m, err
		}
		return n, nil
	}
	return r.reader.Read(p)
}

// detectDelimiter detects the delimiter in a given string.
// It assumes the delimiter is the most frequent non-alphanumeric character.
func detectDelimiter(data string) (rune, error) {
	// Strip BOM from data if present (for delimiter detection)
	if len(data) >= 3 && data[0] == '\xEF' && data[1] == '\xBB' && data[2] == '\xBF' {
		data = data[3:]
	}

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
// Returns the csv.Reader and a cleanup function that must be called to close the file.
// Automatically handles UTF-8 BOM if present.
func CSVReader(filename string, delim ...rune) (*csv.Reader, func() error, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	if !scanner.Scan() {
		file.Close()
		return nil, nil, fmt.Errorf("file is empty")
	}
	firstLine := scanner.Text()

	var delimiter rune
	if len(delim) > 0 {
		delimiter = delim[0]
	} else {
		delimiter, err = detectDelimiter(firstLine)
		if err != nil {
			file.Close()
			return nil, nil, err
		}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		file.Close()
		return nil, nil, fmt.Errorf("could not seek to beginning of file: %w", err)
	}

	// Wrap file in BOM-skipping reader, then CRLF normalizer
	bomReader := &bomSkippingReader{reader: file}
	csvReader := csv.NewReader(crlf.NewReader(bomReader))
	csvReader.LazyQuotes = true
	csvReader.Comma = delimiter

	return csvReader, file.Close, nil
}
