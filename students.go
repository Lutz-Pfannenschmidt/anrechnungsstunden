package main

import (
	"anrechnungsstundenberechner/internal/autocsv"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

// Leistungsdaten
type SemesterStudentData struct {
	Nachname        string
	Vorname         string
	Geburtsdatum    string
	Jahr            int
	Abschnitt       string
	Fach            string
	Fachlehrer      string
	Kursart         string
	Kurs            string
	Note            string
	Abiturfach      string
	Wochenstd       string
	ExterneSchulnr  string
	Zusatzkraft     string
	WochenstdZK     string
	Jahrgang        string
	Jahrgänge       string
	Fehlstd         string
	UnentschFehlstd string
	Mahnung         string
	Sortierung      string
	Mahndatum       string
}

// ReadSemesterStudentDataFromFile reads the CSV file and returns a slice of Student structs.
// It uses the autocsv package to handle different delimiters.
// If no delimiter is provided, it will auto-detect the delimiter.
// If multiple delimiters are provided, it will use the first one.
func ReadSemesterStudentDataFromFile(filename string, delim ...rune) ([]SemesterStudentData, error) {
	file, err := autocsv.CSVReader(filename, delim...)
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %w", err)
	}
	return ReadSemesterStudentDataFromReader(file)
}

// readStudents reads the CSV file and returns a slice of Student structs
func ReadSemesterStudentDataFromReader(reader *csv.Reader) ([]SemesterStudentData, error) {
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}
	if len(records[0]) < 22 {
		return nil, fmt.Errorf("expected at least 22 columns, got %d", len(records[0]))
	}
	if !strings.EqualFold(strings.TrimSpace(records[0][0]), "Nachname") {
		return nil, fmt.Errorf("expected 'Nachname' in the first column, got '%s'", []byte(records[0][0]))
	}

	var students []SemesterStudentData
	for _, record := range records[1:] {
		student := SemesterStudentData{
			Nachname:        strings.TrimSpace(record[0]),
			Vorname:         strings.TrimSpace(record[1]),
			Geburtsdatum:    strings.TrimSpace(record[2]),
			Jahr:            mustAtoi(strings.TrimSpace(record[3])),
			Abschnitt:       strings.TrimSpace(record[4]),
			Fach:            strings.TrimSpace(record[5]),
			Fachlehrer:      strings.TrimSpace(record[6]),
			Kursart:         strings.TrimSpace(record[7]),
			Kurs:            strings.TrimSpace(record[8]),
			Note:            strings.TrimSpace(record[9]),
			Abiturfach:      strings.TrimSpace(record[10]),
			Wochenstd:       strings.TrimSpace(record[11]),
			ExterneSchulnr:  strings.TrimSpace(record[12]),
			Zusatzkraft:     strings.TrimSpace(record[13]),
			WochenstdZK:     strings.TrimSpace(record[14]),
			Jahrgang:        strings.TrimSpace(record[15]),
			Jahrgänge:       strings.TrimSpace(record[16]),
			Fehlstd:         strings.TrimSpace(record[17]),
			UnentschFehlstd: strings.TrimSpace(record[18]),
			Mahnung:         strings.TrimSpace(record[19]),
			Sortierung:      strings.TrimSpace(record[20]),
			Mahndatum:       strings.TrimSpace(record[21]),
		}
		students = append(students, student)
	}

	return students, nil
}

// Teilleistungen
type ExamData struct {
	Nachname     string
	Vorname      string
	Geburtsdatum string
	Jahr         int
	Abschnitt    string
	Fach         string
	Datum        string
	Teilleistung string
	Note         string
	Bemerkung    string
	Lehrkraft    string
}

// ReadExamDataFromFile reads the CSV file and returns a slice of ExamData structs.
// It uses the autocsv package to handle different delimiters.
// If no delimiter is provided, it will auto-detect the delimiter.
// If multiple delimiters are provided, it will use the first one.
func ReadExamDataFromFile(filename string, delim ...rune) ([]ExamData, error) {
	file, err := autocsv.CSVReader(filename, delim...)
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %w", err)
	}
	return ReadExamDataFromReader(file)
}

// ReadExamDataFromReader reads the CSV file and returns a slice of ExamData structs
func ReadExamDataFromReader(reader *csv.Reader) ([]ExamData, error) {
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}
	if len(records[0]) < 11 {
		return nil, fmt.Errorf("expected at least 11 columns, got %d", len(records[0]))
	}
	if !strings.EqualFold(strings.TrimSpace(records[0][0]), "Nachname") {
		return nil, fmt.Errorf("expected 'Nachname' in the first column, got '%s'", []byte(records[0][0]))
	}

	var studentNotes []ExamData
	for _, record := range records[1:] {
		studentNote := ExamData{
			Nachname:     strings.TrimSpace(record[0]),
			Vorname:      strings.TrimSpace(record[1]),
			Geburtsdatum: strings.TrimSpace(record[2]),
			Jahr:         mustAtoi(strings.TrimSpace(record[3])),
			Abschnitt:    strings.TrimSpace(record[4]),
			Fach:         strings.TrimSpace(record[5]),
			Datum:        strings.TrimSpace(record[6]),
			Teilleistung: strings.TrimSpace(record[7]),
			Note:         strings.TrimSpace(record[8]),
			Bemerkung:    strings.TrimSpace(record[9]),
			Lehrkraft:    strings.TrimSpace(record[10]),
		}
		studentNotes = append(studentNotes, studentNote)
	}

	return studentNotes, nil
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
