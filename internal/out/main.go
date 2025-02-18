package out

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/user0608/excel2pdf"
	"github.com/xuri/excelize/v2"
)

type RowData struct {
	Subject  string
	Grade    string
	Students int
	Points   int
}

func RenderTemplate(templatePath, outputPath, yearStr, name string, weekly_hours float64, lead int, lead_points, hours float64, data []RowData) error {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return err
	}

	defer f.Close()

	sheet := f.GetSheetList()[0]
	f.SetCellStr(sheet, "A2", yearStr)
	f.SetCellStr(sheet, "A3", "Name: "+toUpperCaseFirst(name))
	f.SetCellStr(sheet, "D3", floatToString(weekly_hours))

	if len(data) > 8 {
		data = data[:8]
	}

	sum5 := 0.0

	off := 8
	for i, row := range data {
		f.SetCellStr(sheet, "A"+strconv.Itoa(off+i), fmt.Sprintf("%s, %s", toUpperCaseFirst(row.Subject), row.Grade))
		f.SetCellStr(sheet, "B"+strconv.Itoa(off+i), intToString(row.Students))
		f.SetCellStr(sheet, "C"+strconv.Itoa(off+i), intToString(row.Points))
		f.SetCellStr(sheet, "D"+strconv.Itoa(off+i), intToString(row.Points*row.Students))
		sum5 += float64(row.Points * row.Students)
	}

	points := float64(lead) / 100.0 * lead_points
	sum5 += points
	base := weekly_hours * 7.84
	sum7 := max(sum5-base, 0)
	f.SetCellStr(sheet, "B16", fmt.Sprintf("in %% x %s", floatToString(lead_points)))
	f.SetCellStr(sheet, "D16", fmt.Sprintf("%s%% * %s = %s", intToString(lead), floatToString(lead_points), floatToString(points)))
	f.SetCellStr(sheet, "D17", floatToString(sum5))
	f.SetCellStr(sheet, "D18", floatToString(base))
	f.SetCellStr(sheet, "D19", floatToString(sum7))
	f.SetCellStr(sheet, "D20", floatToString(hours))

	tempName := os.TempDir() + "/" + TempName(".xlsx")
	f.SaveAs(tempName)
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

func toUpperCaseFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
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
