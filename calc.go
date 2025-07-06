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
	p, ok := ep[grade][subject]

	if ok {
		return p, ok
	}

	if strings.HasSuffix(subject, " LK") {
		p, ok := ep[grade]["*LK"]

		if ok {
			return p, ok
		}
		return 0, false
	}

	if strings.HasSuffix(subject, " GK") {
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
		return fmt.Errorf("failed to fetch year record: %w", err)
	}
	if year.ID == "" {
		return fmt.Errorf("year record with ID %s not found", yearID)
	}

	td_records := []TeacherDataRecord{}
	err = app.DB().
		Select("teacher_data.id", "semester", "user", "avg_time", "class_lead", "short", "name"). // Include short and name for better usability
		From("teacher_data").
		AndWhere(dbx.NewExp("semester={:semester}", dbx.Params{"semester": yearID})).
		Join("INNER JOIN", "users", dbx.NewExp("users.id = teacher_data.user")).
		All(&td_records)
	if err != nil {
		return fmt.Errorf("failed to fetch teacher data records: %w", err)
	}

	files, err := app.FindRecordsByFilter("files", "year={:year} && semester={:sem}", "", 0, 0, dbx.Params{
		"year": year.StartYear,
		"sem":  year.Semester,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch files for year %s: %w", yearID, err)
	}

	if len(files) != 3 {
		return fmt.Errorf("expected 3 files (exam, course, hours) for year %s, got %d", yearID, len(files))
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
		return fmt.Errorf("failed to read student data (Leistungsdaten) from course file: %w", err)
	}

	if len(semStuData) == 0 {
		return fmt.Errorf("no student data found in course file %s", courseFilePath)
	}

	examData, err := ReadExamDataFromFile(examFilePath)
	if err != nil {
		return fmt.Errorf("failed to read exam data (Teilleistungen) from file %s: %w", examFilePath, err)
	}
	if len(examData) == 0 {
		return fmt.Errorf("no exam data found in file %s", examFilePath)
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
		return fmt.Errorf("failed to fetch exam point records: %w", err)
	}

	if len(examPointRecords) == 0 {
		return fmt.Errorf("no exam point records found")
	}

	for _, record := range examPointRecords {
		examPoints.Set(record.Grade, record.Class, record.Points)
	}

	for _, e := range examData {
		if e.Note == "" || e.Lehrkraft == "" || e.Fach == "" || e.Abschnitt == "" || e.Jahr == year.StartYear || strings.Contains(strings.ToLower(e.Teilleistung), "somi") {
			continue
		}

		parts := strings.Split(e.Nachname, "#")
		if len(parts) != 2 || len(parts[1]) < 2 {
			fmt.Printf("Skipping exam entry with invalid name format: %s\n", e.Nachname)
			continue
		}

		grade := strings.ToUpper(strings.TrimSpace(parts[1])[:2])
		subject := strings.ToUpper(strings.TrimSpace(e.Fach))

		mul, _ := examPoints.Get(grade, subject)

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
				data[uID][i].Exams++
				found = true
				break
			}
		}
		if !found {
			data[uID] = append(data[uID], render.RowData{Subject: subject, Grade: grade, Exams: 1, Points: mul})
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
		return fmt.Errorf("failed to remove output directory %s: %w", outDir, err)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outDir, err)
	}

	for id, td := range tdByID {
		fmt.Printf("Rendering PDF for user %s (%s, %s)\n", td.Name, td.ID, id)
		_, ok := scores[id]
		if !ok {
			fmt.Printf("No score for user %s (%s), skipping rendering\n", td.Name, id)
			continue
		}
	}

	for uID, score := range scores {
		td, ok := tdByID[uID]
		if !ok {
			return fmt.Errorf("teacher data for user %s not found", uID)
		}
		d, ok := data[uID]
		if !ok {
			return fmt.Errorf("data for user %s not found", uID)
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
			fmt.Printf("Failed to render template for user %s (%s): %v\n", nameParts[0]+" ("+td.UserID+")", uID, err)
		}
	}

	fmt.Printf("Rendered %d PDFs for year %s\n", len(scores), yearID)

	// Combine all PDFs into a single PDF
	combinedPDFPath := path.Join(outDir, "combined.pdf")
	if err := render.CombinePDFsFromDir(outDir, combinedPDFPath); err != nil {
		return fmt.Errorf("failed to combine PDFs: %w", err)
	}

	// Save the individual PDFs to the database
	pdfsCollection, err := app.FindCollectionByNameOrId("pdfs")
	if err != nil {
		return fmt.Errorf("failed to find pdfs collection: %w", err)
	}

	for uID := range scores {
		pdfPath := path.Join(outDir, uID+".pdf")
		pdfFile, err := filesystem.NewFileFromPath(pdfPath)
		if err != nil {
			return fmt.Errorf("failed to create file from path %s: %w", pdfPath, err)
		}

		record := core.NewRecord(pdfsCollection)
		record.Set("user", uID)
		record.Set("semester", yearID)
		record.Set("pdf", pdfFile)

		if err := app.Save(record); err != nil {
			return fmt.Errorf("failed to save PDF record for user %s: %w", uID, err)
		}
	}

	// Save the results record
	resultsCollection, err := app.FindCollectionByNameOrId("results")
	if err != nil {
		return fmt.Errorf("failed to find results collection: %w", err)
	}

	pdf, err := filesystem.NewFileFromPath(combinedPDFPath)
	if err != nil {
		return fmt.Errorf("failed to create file from path %s: %w", combinedPDFPath, err)
	}

	humanReadableScores := make(map[string]float64)
	untisScores := make(map[string]float64)
	for uID, score := range scores {
		td, ok := tdByID[uID]
		if !ok {
			return fmt.Errorf("teacher data for user %s not found", uID)
		}

		nameParts := strings.Split(td.Name, "_NAME_COLLISION_")
		name := nameParts[0] + " (" + strings.ToUpper(td.Short) + ")"

		humanReadableScores[name] = score
		untisScores[strings.ToUpper(td.Short)] = score
	}

	untisPath := path.Join(outDir, "untis.txt")
	if err := render.WriteToUntisDataFile(untisScores, untisPath); err != nil {
		return fmt.Errorf("failed to write Untis data file: %w", err)
	}

	untisFile, err := filesystem.NewFileFromPath(untisPath)
	if err != nil {
		return fmt.Errorf("failed to create file from path %s: %w", untisPath, err)
	}

	resRecord := core.NewRecord(resultsCollection)
	resRecord.Set("semester", yearID)
	resRecord.Set("data", humanReadableScores)
	resRecord.Set("pdf", pdf)
	resRecord.Set("untis", untisFile)

	if err := app.Save(resRecord); err != nil {
		return fmt.Errorf("failed to save results record: %w", err)
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
