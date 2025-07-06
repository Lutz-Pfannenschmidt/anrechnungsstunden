package render

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/user0608/excel2pdf"
	"github.com/xuri/excelize/v2"
)

type RowData struct {
	Subject string
	Grade   string
	Exams   int
	Points  float64
}

func RenderTemplate(templateFile []byte, outputPath, yearStr, name string, weekly_hours, baseMultiplier, extraPoints, leadPercentage, pointsPerLead, hours float64, data []RowData) error {
	reader := strings.NewReader(string(templateFile))

	f, err := excelize.OpenReader(reader)
	if err != nil {
		return err
	}

	defer f.Close()

	sheet := f.GetSheetList()[0]
	f.SetCellStr(sheet, "A2", yearStr)
	f.SetCellStr(sheet, "A3", "Name: "+name)
	f.SetCellStr(sheet, "D3", floatToString(weekly_hours))

	if len(data) > 15 {
		data = data[:15]
	}

	sum5 := 0.0

	off := 8
	for i, row := range data {
		f.SetCellStr(sheet, "A"+strconv.Itoa(off+i), fmt.Sprintf("%s, %s", row.Subject, row.Grade))
		f.SetCellStr(sheet, "B"+strconv.Itoa(off+i), intToString(row.Exams))
		f.SetCellStr(sheet, "C"+strconv.Itoa(off+i), floatToString(row.Points))
		f.SetCellStr(sheet, "D"+strconv.Itoa(off+i), floatToString(row.Points*float64(row.Exams)))
		sum5 += row.Points * float64(row.Exams)
	}

	points := float64(leadPercentage) / 100.0 * pointsPerLead
	sum5 += points + extraPoints
	base := weekly_hours * baseMultiplier
	sum7 := max(sum5-base, 0)
	f.SetCellStr(sheet, "D23", floatToString(extraPoints))
	f.SetCellStr(sheet, "C24", fmt.Sprintf("in %% x %s", floatToString(pointsPerLead)))
	f.SetCellStr(sheet, "D24", fmt.Sprintf("%s%% * %s = %s", floatToString(leadPercentage), floatToString(pointsPerLead), floatToString(points)))
	f.SetCellStr(sheet, "D25", floatToString(sum5))
	f.SetCellStr(sheet, "C26", fmt.Sprintf("Erteilte Stunden x %s", floatToString(baseMultiplier)))
	f.SetCellStr(sheet, "D26", floatToString(base))
	f.SetCellStr(sheet, "D27", floatToString(sum7))
	f.SetCellStr(sheet, "D28", floatToString(hours))

	tempName := os.TempDir() + "/" + TempName(".xlsx")
	err = f.SaveAs(tempName)
	if err != nil {
		return fmt.Errorf("failed to save temporary Excel file: %w", err)
	}
	pdfFilePath, err := excel2pdf.ConvertExcelToPdf(tempName)
	if err != nil {
		return err
	}
	err = os.Remove(tempName)
	if err != nil {
		return err
	}

	err = os.Rename(pdfFilePath, outputPath)
	if err != nil {
		return err
	}

	return nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func TempName(extension string) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 8

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b) + extension
}

func floatToString(englishNumber float64) string {
	return stringToString(strconv.FormatFloat(englishNumber, 'f', 3, 64))
}

func intToString(englishNumber int) string {
	return stringToString(strconv.Itoa(englishNumber))
}

func stringToString(englishNumber string) string {
	parts := strings.Split(englishNumber, ".")
	whole := strings.ReplaceAll(parts[0], ",", ".")
	decimal := ""
	if len(parts) > 1 {
		decimal = strings.ReplaceAll(parts[1], ",", "")
		decimal = strings.TrimRight(decimal, "0")
	}
	if decimal == "" {
		return whole
	} else {
		return fmt.Sprintf("%s,%s", whole, decimal)
	}
}

func WriteToUntisDataFile(scores map[string]float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating Untis data file: %w", err)
	}
	defer file.Close()

	i := 0
	for teacher, hours := range scores {
		line := fmt.Sprintf("%d;;;;\"%s\";\"500\",\"%.3f;;;;\"LK\";0.000;0.000;;%.3f\n", i+10000, strings.ToUpper(teacher), hours, hours)
		i++
		_, err := file.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing to Untis data file: %w", err)
		}
	}

	return nil
}

func CombinePDFs(inputPaths []string, outputPath string) error {
	config := model.NewDefaultConfiguration()
	err := api.MergeCreateFile(inputPaths, outputPath, false, config)
	if err != nil {
		return fmt.Errorf("error merging PDFs: %w", err)
	}
	return nil
}

func CombinePDFsFromDir(dir string, outputPath string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", dir, err)
	}
	var inputPaths []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			inputPaths = append(inputPaths, path.Join(dir, file.Name()))
		}
	}
	if len(inputPaths) == 0 {
		return fmt.Errorf("no PDF files found in directory %s", dir)
	}
	return CombinePDFs(inputPaths, outputPath)
}
