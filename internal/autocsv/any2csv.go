package autocsv

import (
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

// stripBOM removes UTF-8 BOM from the beginning of a string if present
func stripBOM(s string) string {
	if len(s) >= 3 && s[0] == '\xEF' && s[1] == '\xBB' && s[2] == '\xBF' {
		return s[3:]
	}
	return s
}

// normalizeCRLF converts Windows line endings (CRLF) to Unix line endings (LF)
func normalizeCRLF(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// escapeCSVField properly escapes a CSV field value.
// Fields containing commas, quotes, or newlines are quoted,
// and internal quotes are doubled.
func escapeCSVField(field string) string {
	needsQuoting := strings.ContainsAny(field, ",\"\n\r")
	if needsQuoting {
		// Double any existing quotes and wrap in quotes
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return field
}

func ReadXLSXFileToCSV(path string) (string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return "", fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return "", err
	}

	var result strings.Builder

	for _, row := range rows {
		escapedFields := make([]string, len(row))
		for i, field := range row {
			escapedFields[i] = escapeCSVField(field)
		}
		result.WriteString(strings.Join(escapedFields, ","))
		result.WriteString("\n")
	}

	return result.String(), nil
}

func ReadCSVFile(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	content := string(file)
	// Handle BOM and CRLF
	content = stripBOM(content)
	content = normalizeCRLF(content)
	return content, nil
}

func ReadAnyFileToCSV(path string) (string, error) {
	if strings.HasSuffix(strings.ToLower(path), ".csv") {
		return ReadCSVFile(path)
	} else if strings.HasSuffix(strings.ToLower(path), ".xlsx") || strings.HasSuffix(strings.ToLower(path), ".xls") {
		return ReadXLSXFileToCSV(path)
	} else {
		return "", fmt.Errorf("Nicht unterstützter Dateityp: %s. Erlaubt sind .csv, .xlsx und .xls", path)
	}
}
