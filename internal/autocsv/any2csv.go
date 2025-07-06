package autocsv

import (
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

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

	result := ""

	for _, row := range rows {
		result += strings.Join(row, ",") + "\n"
	}

	return result, nil
}

func ReadCSVFile(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func ReadAnyFileToCSV(path string) (string, error) {
	if strings.HasSuffix(path, ".csv") {
		return ReadCSVFile(path)
	} else if strings.HasSuffix(path, ".xlsx") || strings.HasSuffix(path, ".xls") {
		return ReadXLSXFileToCSV(path)
	} else {
		return "", fmt.Errorf("unsupported file type")
	}
}
