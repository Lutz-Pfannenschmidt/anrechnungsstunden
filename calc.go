package main

import (
	"anrechnungsstundenberechner/internal/render"
	_ "embed"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

//go:embed templates/template.xlsx
var templateFile []byte

type GradeToSubjectToPoints map[string]map[string]float64

func (ep GradeToSubjectToPoints) Set(grade, subject string, points float64) {
	if _, ok := ep[grade]; !ok {
		ep[grade] = map[string]float64{}
	}
	ep[grade][subject] = points
}

func (ep GradeToSubjectToPoints) Get(grade, subject string) (float64, bool) {
	grade = strings.ToUpper(grade)
	subject = strings.ToUpper(subject)

	p, ok := ep[grade][subject]

	if ok {
		return p, ok
	}

	if strings.HasSuffix((subject), " LK") {
		p, ok := ep[grade]["*LK"]

		if ok {
			return p, ok
		}
		return 0, false
	}

	if strings.HasSuffix((subject), " GK") {
		p, ok := ep[grade]["*GK"]

		if ok {
			return p, ok
		}
		return 0, false
	}

	return 0, false
}

// This function:
//   - calculates the results for a given year and semester (distributes the hours)
//   - renders the results into a PDF for each user
//   - combines the results into a single PDF for the year
//   - saves all the above into the database
func CalculateAndRenderResults(app core.App, yearID string) error {

	year := YearsRecord{}
	err := app.DB().
		Select("*").
		From("years").
		AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": yearID})).
		One(&year)
	if err != nil {
		return fmt.Errorf("Halbjahresabschluss mit ID '%s' konnte nicht geladen werden: %w", yearID, err)
	}
	if year.ID == "" {
		return fmt.Errorf("Halbjahresabschluss mit ID '%s' existiert nicht in der Datenbank", yearID)
	}

	td_records := []TeacherDataRecord{}
	err = app.DB().
		Select("teacher_data.id", "semester", "user", "avg_time", "class_lead", "add_points", "short", "name"). // Include short and name for better usability
		From("teacher_data").
		AndWhere(dbx.NewExp("semester={:semester}", dbx.Params{"semester": yearID})).
		Join("INNER JOIN", "users", dbx.NewExp("users.id = teacher_data.user")).
		All(&td_records)
	if err != nil {
		return fmt.Errorf("Lehrerdaten für Halbjahr '%s' konnten nicht geladen werden: %w", yearID, err)
	}

	if len(td_records) == 0 {
		return fmt.Errorf("Keine Lehrerdaten für Halbjahr '%s' gefunden. Bitte wählen Sie zuerst Lehrkräfte aus (Schritt 2)", yearID)
	}

	files, err := app.FindRecordsByFilter("files", "year={:year} && semester={:sem}", "", 0, 0, dbx.Params{
		"year": year.StartYear,
		"sem":  year.Semester,
	})
	if err != nil {
		return fmt.Errorf("Dateien für Schuljahr %d/%d, %d. Halbjahr konnten nicht geladen werden: %w", year.StartYear, year.StartYear+1, year.Semester, err)
	}

	if len(files) != 3 {
		foundTypes := []string{}
		for _, f := range files {
			foundTypes = append(foundTypes, f.GetString("type"))
		}
		return fmt.Errorf("Es werden 3 Dateien (exam, course, hours) für Schuljahr %d/%d benötigt, aber nur %d gefunden: %v. Bitte laden Sie alle erforderlichen Dateien hoch",
			year.StartYear, year.StartYear+1, len(files), foundTypes)
	}

	var examFilePath, courseFilePath string

	for _, file := range files {
		if file.GetString("type") == string(FilesRecordTypeExam) {
			examFilePath = app.DataDir() + "/storage/" + file.BaseFilesPath() + "/" + file.GetString("file")
		} else if file.GetString("type") == string(FilesRecordTypeCourse) {
			courseFilePath = app.DataDir() + "/storage/" + file.BaseFilesPath() + "/" + file.GetString("file")
		}
	}

	semStuData, err := ReadSemesterStudentDataFromFile(courseFilePath)
	if err != nil {
		return fmt.Errorf("Fehler beim Lesen der Leistungsdaten (Kursdaten) aus '%s': %w. Bitte überprüfen Sie das Dateiformat", courseFilePath, err)
	}

	if len(semStuData) == 0 {
		return fmt.Errorf("Keine Schülerdaten in der Leistungsdaten-Datei '%s' gefunden. Die Datei scheint leer oder ungültig zu sein", courseFilePath)
	}

	examData, err := ReadExamDataFromFile(examFilePath)
	fmt.Printf("Reading exam data from file %s\n", examFilePath)
	if err != nil {
		return fmt.Errorf("Fehler beim Lesen der Teilleistungen (Klausurdaten) aus '%s': %w. Bitte überprüfen Sie das Dateiformat", examFilePath, err)
	}
	if len(examData) == 0 {
		return fmt.Errorf("Keine Klausurdaten in der Teilleistungen-Datei '%s' gefunden. Die Datei scheint leer oder ungültig zu sein", examFilePath)
	}

	scores := map[string]float64{}           // UserID -> Score
	data := map[string][]render.RowData{}    // UserID -> RowData
	uIDs := map[string]string{}              // Short (uppercase) -> UserID
	shorts := map[string]string{}            // UID -> Short (uppercase)
	tdByID := map[string]TeacherDataRecord{} // UserID -> TeacherDataRecord

	for _, td := range td_records {
		scores[td.UserID] = 0
		data[td.UserID] = []render.RowData{}
		uIDs[strings.ToUpper(td.Short)] = td.UserID
		tdByID[td.UserID] = td
		shorts[td.UserID] = strings.ToUpper(td.Short)
	}

	examPoints := GradeToSubjectToPoints{}

	examPointRecords := []PartPointsRecord{}
	err = app.DB().
		Select("*").
		From("part_points").
		AndWhere(dbx.NewExp("points>0")).
		All(&examPointRecords)
	if err != nil {
		return fmt.Errorf("Fehler beim Laden der Punktekonfiguration (part_points): %w", err)
	}

	if len(examPointRecords) == 0 {
		return fmt.Errorf("Keine Punktekonfiguration gefunden. Bitte konfigurieren Sie zuerst die Punkte pro Klausur unter 'Punktekonfiguration'")
	}

	for _, record := range examPointRecords {
		examPoints.Set(record.Grade, record.Class, record.Points)
	}

	for _, e := range examData {
		if e.Note == "" || e.Lehrkraft == "" || e.Fach == "" || e.Abschnitt == "" || e.Jahr != year.StartYear || strings.Contains(strings.ToLower(e.Teilleistung), "somi") {
			continue
		}

		parts := strings.Split(e.Nachname, "#")
		if len(parts) != 2 || len(parts[1]) < 2 {
			fmt.Printf("Skipping exam entry with invalid name format: %s\n", e.Nachname)
			continue
		}

		fmt.Printf("Processing exam entry for %s %s (%s) in %s\n", parts[0], parts[1], e.Fach, e.Abschnitt)

		grade := strings.ToUpper(strings.TrimSpace(parts[1])[:2])
		subject := strings.ToUpper(strings.TrimSpace(e.Fach))

		mul, _ := examPoints.Get(grade, subject)
		examsCount := 1

		if grade == "EF" || grade == "Q1" || (grade == "Q2" && year.Semester != 2) {
			for _, s := range semStuData {
				if GermanEquals(s.Nachname, e.Nachname) && GermanEquals(s.Vorname, e.Vorname) && GermanEquals(s.Geburtsdatum, e.Geburtsdatum) && strings.EqualFold(e.Fach, s.Fach) && strings.EqualFold(e.Lehrkraft, s.Fachlehrer) {
					ka := strings.TrimSpace(s.Kursart)
					if strings.EqualFold(ka, "LK1") || strings.EqualFold(ka, "LK2") {
						mul, _ = examPoints.Get(grade, subject+" LK")
					} else if strings.EqualFold(ka, "GKS") {
						mul, _ = examPoints.Get(grade, subject+" GK")
					} else if strings.EqualFold(ka, "GKM") {
						mul = 0
					}
					break
				}
			}
		} else if grade == "Q2" && year.Semester == 2 {
			for _, s := range semStuData {
				if GermanEquals(s.Nachname, e.Nachname) && GermanEquals(s.Vorname, e.Vorname) && GermanEquals(s.Geburtsdatum, e.Geburtsdatum) && strings.EqualFold(e.Fach, s.Fach) && strings.EqualFold(e.Lehrkraft, s.Fachlehrer) {
					ka := strings.TrimSpace(s.Kursart)
					if strings.EqualFold(ka, "LK1") || strings.EqualFold(ka, "LK2") {
						mul, _ = examPoints.Get(grade, subject+" LK")
					} else if strings.EqualFold(ka, "AB3") {
						mul, _ = examPoints.Get(grade, subject+" GK")
					} else {
						mul = 0
					}
					break
				}
			}
			examsCount = 2 // Q2 in the second semester counts as two exams
		}

		short := strings.ToUpper(strings.TrimSpace(e.Lehrkraft))

		uID, ok := uIDs[short]
		if !ok {
			continue
		}

		if _, ok := scores[uID]; !ok {
			continue
		}

		found := false
		for i, d := range data[uID] {
			if strings.EqualFold(grade, d.Grade) && strings.EqualFold(subject, d.Subject) && d.Points == mul {
				data[uID][i].Exams += examsCount
				found = true
				break
			}
		}
		if !found {
			data[uID] = append(data[uID], render.RowData{Subject: subject, Grade: grade, Exams: examsCount, Points: mul})
		}

		scores[uID] += mul
	}

	for scoresUID := range scores {

		td, ok := tdByID[scoresUID]
		if !ok {
			continue
		}

		scores[scoresUID] += td.ClassLead/100*year.LeadPoints + td.AddPoints - td.AverageTime*year.BaseMul
		if scores[scoresUID] < 0 {
			scores[scoresUID] = 0
		}
	}

	fmt.Printf("Calculated scores for %d teachers in year %s\n", len(scores), yearID)

	for uID, score := range scores {
		fmt.Printf("Score for user %s: %f\n", shorts[uID], score)
	}

	scores = distribute(year.TotalPoints, scores)
	fmt.Printf("Distributed scores for %d teachers in year %s\n", len(scores), yearID)

	for uID, score := range scores {
		fmt.Printf("Score for user %s: %f\n", shorts[uID], score)
	}

	semesterStr := fmt.Sprintf("%d. Halbjahr", year.Semester)
	yearStr := fmt.Sprintf("%d/%s, %s", year.StartYear, strconv.Itoa(year.StartYear + 1)[2:], semesterStr)

	outDir := path.Join(os.TempDir(), "Anrechnungsstundenberechner", year.ID+"_"+strconv.Itoa(year.Semester))

	if err := os.RemoveAll(outDir); err != nil {
		return fmt.Errorf("Temporäres Ausgabeverzeichnis '%s' konnte nicht gelöscht werden: %w", outDir, err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("Temporäres Ausgabeverzeichnis '%s' konnte nicht erstellt werden: %w", outDir, err)
	}

	for id, td := range tdByID {
		fmt.Printf("Rendering PDF for user %s (%s, %s)\n", td.Name, td.ID, id)
		_, ok := scores[id]
		if !ok {
			fmt.Printf("No score for user %s (%s), skipping rendering\n", td.Name, id)
			continue
		}
	}

	var renderErrors []string
	for uID, score := range scores {
		td, ok := tdByID[uID]
		if !ok {
			return fmt.Errorf("Lehrerdaten für Benutzer-ID '%s' nicht gefunden. Datenbankinkonsistenz erkannt", uID)
		}
		d, ok := data[uID]
		if !ok {
			return fmt.Errorf("Klausurdaten für Benutzer '%s' (%s) nicht gefunden", td.Name, uID)
		}

		nameParts := strings.Split(td.Name, "_NAME_COLLISION_")
		name := nameParts[0] + " (" + strings.ToUpper(td.Short) + ")"

		err := render.RenderTemplate(
			templateFile,
			path.Join(outDir, uID+".pdf"),
			yearStr,
			name,
			td.AverageTime,
			year.BaseMul,
			td.AddPoints,
			td.ClassLead,
			year.LeadPoints,
			score,
			d,
		)
		if err != nil {
			renderErrors = append(renderErrors, fmt.Sprintf("- %s (%s): %v", nameParts[0], td.Short, err))
		}
	}

	if len(renderErrors) > 0 {
		return fmt.Errorf("PDF-Erstellung fehlgeschlagen für %d Lehrkraft/Lehrkräfte:\n%s", len(renderErrors), strings.Join(renderErrors, "\n"))
	}

	fmt.Printf("Rendered %d PDFs for year %s\n", len(scores), yearID)

	// Combine all PDFs into a single PDF
	combinedPDFPath := path.Join(outDir, "combined.pdf")
	if err := render.CombinePDFsFromDir(outDir, combinedPDFPath); err != nil {
		return fmt.Errorf("Fehler beim Zusammenführen der PDFs zu einer Gesamtdatei: %w", err)
	}

	// Save the individual PDFs to the database
	pdfsCollection, err := app.FindCollectionByNameOrId("pdfs")
	if err != nil {
		return fmt.Errorf("Datenbankfehler: 'pdfs' Collection nicht gefunden: %w", err)
	}

	for uID := range scores {
		pdfPath := path.Join(outDir, uID+".pdf")
		pdfFile, err := filesystem.NewFileFromPath(pdfPath)
		if err != nil {
			return fmt.Errorf("PDF-Datei '%s' konnte nicht gelesen werden: %w", pdfPath, err)
		}

		record := core.NewRecord(pdfsCollection)
		record.Set("user", uID)
		record.Set("semester", yearID)
		record.Set("pdf", pdfFile)

		if err := app.Save(record); err != nil {
			td := tdByID[uID]
			return fmt.Errorf("PDF für Lehrkraft '%s' (%s) konnte nicht gespeichert werden: %w", td.Name, td.Short, err)
		}
	}

	// Save the results record
	resultsCollection, err := app.FindCollectionByNameOrId("results")
	if err != nil {
		return fmt.Errorf("Datenbankfehler: 'results' Collection nicht gefunden: %w", err)
	}

	pdf, err := filesystem.NewFileFromPath(combinedPDFPath)
	if err != nil {
		return fmt.Errorf("Kombinierte PDF-Datei '%s' konnte nicht gelesen werden: %w", combinedPDFPath, err)
	}

	humanReadableScores := make(map[string]float64)
	untisScores := make(map[string]float64)
	for uID, score := range scores {
		td, ok := tdByID[uID]
		if !ok {
			return fmt.Errorf("Lehrerdaten für Benutzer-ID '%s' bei Ergebniserstellung nicht gefunden", uID)
		}

		nameParts := strings.Split(td.Name, "_NAME_COLLISION_")
		name := nameParts[0] + " (" + strings.ToUpper(td.Short) + ")"

		humanReadableScores[name] = score
		untisScores[strings.ToUpper(td.Short)] = score
	}

	untisPath := path.Join(outDir, "untis.txt")
	if err := render.WriteToUntisDataFile(untisScores, untisPath); err != nil {
		return fmt.Errorf("Fehler beim Erstellen der Untis-Export-Datei: %w", err)
	}

	untisFile, err := filesystem.NewFileFromPath(untisPath)
	if err != nil {
		return fmt.Errorf("Untis-Export-Datei '%s' konnte nicht gelesen werden: %w", untisPath, err)
	}

	resRecord := core.NewRecord(resultsCollection)
	resRecord.Set("semester", yearID)
	resRecord.Set("data", humanReadableScores)
	resRecord.Set("pdf", pdf)
	resRecord.Set("untis", untisFile)

	if err := app.Save(resRecord); err != nil {
		return fmt.Errorf("Ergebnisdatensatz konnte nicht gespeichert werden: %w", err)
	}

	return nil
}

func distribute(points float64, teacherScores map[string]float64) map[string]float64 {
	res := make(map[string]float64)
	var d float64
	for teacher, score := range teacherScores {
		d += score
		res[teacher] = score
	}

	if d == 0 {
		fmt.Println("No scores to distribute, returning empty map")
		return res
	}

	for teacher, score := range res {
		res[teacher] = score * points / d
	}
	return res
}
