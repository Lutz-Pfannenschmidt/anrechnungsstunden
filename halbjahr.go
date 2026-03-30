package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// CreateHalbjahrRequest represents the request body for creating a new Halbjahresabschluss
type CreateHalbjahrRequest struct {
	Year      int    `json:"year"`
	Semester  int    `json:"semester"`
	SplitDate string `json:"split_date"` // Format: DD.MM.YYYY
}

// SelectTeachersRequest represents the request body for selecting teachers
type SelectTeachersRequest struct {
	YearID   string                      `json:"year_id"`
	Teachers map[string]TeacherSelection `json:"teachers"` // Name -> TeacherSelection
}

type TeacherSelection struct {
	Short   string  `json:"short"`
	Email   string  `json:"email"`
	AvgTime float64 `json:"avg_time"`
}

// CloseYearRequest represents the request body for closing a year
type CloseYearRequest struct {
	YearID      string                      `json:"year_id"`
	TeacherData map[string]TeacherCloseData `json:"teacher_data"` // teacher_data ID -> data
	BaseMul     float64                     `json:"base_mul"`
	LeadPoints  float64                     `json:"lead_points"`
	TotalPoints float64                     `json:"total_points"`
}

type TeacherCloseData struct {
	ClassLead float64 `json:"class_lead"`
	AddPoints float64 `json:"add_points"`
}

// createHalbjahr handles the creation of a new Halbjahresabschluss
// This endpoint expects multipart/form-data with files and JSON data
func createHalbjahr(e *core.RequestEvent) error {
	// Parse form data
	if err := e.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		return apis.NewBadRequestError("Fehler beim Parsen der Formulardaten: "+err.Error(), nil)
	}

	// Get form values
	yearStr := e.Request.FormValue("year")
	semesterStr := e.Request.FormValue("semester")
	splitDateStr := e.Request.FormValue("split_date")

	// Validate year
	var year int
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil || year < 2000 {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiges Jahr '%s'. Das Jahr muss eine Zahl größer als 2000 sein.", yearStr),
			map[string]any{"field": "year", "value": yearStr},
		)
	}

	// Validate semester
	var semester int
	if _, err := fmt.Sscanf(semesterStr, "%d", &semester); err != nil || semester < 1 || semester > 2 {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiges Halbjahr '%s'. Das Halbjahr muss 1 oder 2 sein.", semesterStr),
			map[string]any{"field": "semester", "value": semesterStr},
		)
	}

	// Parse and validate split date
	splitDate, err := parseGermanDate(splitDateStr)
	if err != nil {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiges Datum '%s'. Bitte verwenden Sie das Format TT.MM.JJJJ (z.B. 09.02.2026).", splitDateStr),
			map[string]any{"field": "split_date", "value": splitDateStr, "error": err.Error()},
		)
	}

	// Check if date is a Monday
	if splitDate.Weekday() != time.Monday {
		return apis.NewBadRequestError(
			fmt.Sprintf("Das Datum %s ist kein Montag (%s). Bitte wählen Sie den ersten Montag des 2. Halbjahres.",
				splitDateStr, germanWeekday(splitDate.Weekday())),
			map[string]any{"field": "split_date", "value": splitDateStr, "weekday": splitDate.Weekday().String()},
		)
	}

	// Check if date is within the school year
	if splitDate.Year() < year || splitDate.Year() > year+1 {
		return apis.NewBadRequestError(
			fmt.Sprintf("Das Datum %s liegt nicht im Schuljahr %d/%d.", splitDateStr, year, year+1),
			map[string]any{"field": "split_date", "value": splitDateStr, "year": year},
		)
	}

	// Check if this year/semester combination already exists
	existingYear := YearsRecord{}
	err = e.App.DB().
		Select("id").
		From("years").
		Where(dbx.NewExp("start_year={:year} AND semester={:sem}", dbx.Params{"year": year, "sem": semester})).
		One(&existingYear)
	if err == nil && existingYear.ID != "" {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ein Halbjahresabschluss für %d/%d, %d. Halbjahr existiert bereits (ID: %s). Bitte löschen Sie diesen zuerst.", year, year+1, semester, existingYear.ID),
			map[string]any{"existing_id": existingYear.ID},
		)
	}

	// Check for and clean up orphaned file records from previous failed attempts
	// This can happen if a previous createHalbjahr failed midway
	orphanedFiles := []struct {
		ID   string `db:"id"`
		Type string `db:"type"`
	}{}
	err = e.App.DB().
		Select("id", "type").
		From("files").
		Where(dbx.NewExp("year={:year} AND semester={:sem}", dbx.Params{"year": year, "sem": semester})).
		All(&orphanedFiles)
	if err == nil && len(orphanedFiles) > 0 {
		// Delete orphaned files
		orphanedTypes := []string{}
		for _, f := range orphanedFiles {
			orphanedTypes = append(orphanedTypes, f.Type)
			_, delErr := e.App.DB().
				NewQuery("DELETE FROM files WHERE id = {:id}").
				Bind(dbx.Params{"id": f.ID}).
				Execute()
			if delErr != nil {
				return apis.NewApiError(500,
					fmt.Sprintf("Fehler beim Bereinigen verwaister Datei-Einträge (ID: %s): %v", f.ID, delErr),
					map[string]any{"orphaned_file_id": f.ID},
				)
			}
		}
		// Log that we cleaned up orphaned files
		fmt.Printf("Cleaned up %d orphaned file records for year %d semester %d: %v\n", len(orphanedFiles), year, semester, orphanedTypes)
	}

	// Get uploaded files
	entlastungFile, entlastungHeader, err := e.Request.FormFile("file_entlastung")
	if err != nil {
		return apis.NewBadRequestError(
			"Datei 'Entlastungsstunden' fehlt oder konnte nicht gelesen werden: "+err.Error(),
			map[string]any{"field": "file_entlastung"},
		)
	}
	defer entlastungFile.Close()

	leistungFile, leistungHeader, err := e.Request.FormFile("file_leistung")
	if err != nil {
		return apis.NewBadRequestError(
			"Datei 'Leistungsdaten' fehlt oder konnte nicht gelesen werden: "+err.Error(),
			map[string]any{"field": "file_leistung"},
		)
	}
	defer leistungFile.Close()

	teilleistungFile, teilleistungHeader, err := e.Request.FormFile("file_teilleistung")
	if err != nil {
		return apis.NewBadRequestError(
			"Datei 'Teilleistungen' fehlt oder konnte nicht gelesen werden: "+err.Error(),
			map[string]any{"field": "file_teilleistung"},
		)
	}
	defer teilleistungFile.Close()

	// Validate file types
	allowedHoursExts := []string{".csv", ".xls", ".xlsx"}
	allowedDataExts := []string{".csv", ".dat"}

	if !hasValidExtension(entlastungHeader.Filename, allowedHoursExts) {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiger Dateityp für 'Entlastungsstunden': %s. Erlaubt sind: %s",
				entlastungHeader.Filename, strings.Join(allowedHoursExts, ", ")),
			map[string]any{"field": "file_entlastung", "filename": entlastungHeader.Filename},
		)
	}

	if !hasValidExtension(leistungHeader.Filename, allowedDataExts) {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiger Dateityp für 'Leistungsdaten': %s. Erlaubt sind: %s",
				leistungHeader.Filename, strings.Join(allowedDataExts, ", ")),
			map[string]any{"field": "file_leistung", "filename": leistungHeader.Filename},
		)
	}

	if !hasValidExtension(teilleistungHeader.Filename, allowedDataExts) {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiger Dateityp für 'Teilleistungen': %s. Erlaubt sind: %s",
				teilleistungHeader.Filename, strings.Join(allowedDataExts, ", ")),
			map[string]any{"field": "file_teilleistung", "filename": teilleistungHeader.Filename},
		)
	}

	// Create the year record
	yearsCollection, err := e.App.FindCollectionByNameOrId("years")
	if err != nil {
		return apis.NewApiError(500, "Interne Datenbankfehler: 'years' Collection nicht gefunden", map[string]any{"error": err.Error()})
	}

	yearRecord := core.NewRecord(yearsCollection)
	yearRecord.Set("start_year", year)
	yearRecord.Set("semester", semester)
	yearRecord.Set("state", "uploaded")
	yearRecord.Set("split", splitDate)

	if err := e.App.Save(yearRecord); err != nil {
		return apis.NewBadRequestError(
			"Fehler beim Erstellen des Halbjahresabschlusses: "+err.Error(),
			map[string]any{"year": year, "semester": semester},
		)
	}

	// Create file records
	filesCollection, err := e.App.FindCollectionByNameOrId("files")
	if err != nil {
		// Rollback: delete the year record
		_ = e.App.Delete(yearRecord)
		return apis.NewApiError(500, "Interne Datenbankfehler: 'files' Collection nicht gefunden", map[string]any{"error": err.Error()})
	}

	// Track created file record IDs for rollback
	createdFileIDs := []string{}

	// Rollback helper - deletes year record and all created file records
	rollback := func() {
		for _, fileID := range createdFileIDs {
			_, _ = e.App.DB().
				NewQuery("DELETE FROM files WHERE id = {:id}").
				Bind(dbx.Params{"id": fileID}).
				Execute()
		}
		_ = e.App.Delete(yearRecord)
	}

	// Helper to create file record from multipart file
	createFileRecord := func(fileType, formFieldName string) error {
		record := core.NewRecord(filesCollection)
		record.Set("year", year)
		record.Set("semester", semester)
		record.Set("type", fileType)

		// Get the file headers from the multipart form
		fileHeaders := e.Request.MultipartForm.File[formFieldName]
		if len(fileHeaders) == 0 {
			return fmt.Errorf("Keine Datei für '%s' gefunden", formFieldName)
		}

		// Convert multipart.FileHeader to filesystem.File
		files := make([]*filesystem.File, 0, len(fileHeaders))
		for _, fh := range fileHeaders {
			f, err := filesystem.NewFileFromMultipart(fh)
			if err != nil {
				return fmt.Errorf("Fehler beim Verarbeiten der Datei '%s': %w", fh.Filename, err)
			}
			files = append(files, f)
		}

		record.Set("file", files)

		if err := e.App.Save(record); err != nil {
			return fmt.Errorf("Fehler beim Speichern der %s-Datei: %w", fileType, err)
		}

		// Track for rollback
		createdFileIDs = append(createdFileIDs, record.Id)
		return nil
	}

	// Create file records with proper rollback
	if err := createFileRecord("hours", "file_entlastung"); err != nil {
		rollback()
		return apis.NewBadRequestError(err.Error(), nil)
	}

	if err := createFileRecord("course", "file_leistung"); err != nil {
		rollback()
		return apis.NewBadRequestError(err.Error(), nil)
	}

	if err := createFileRecord("exam", "file_teilleistung"); err != nil {
		rollback()
		return apis.NewBadRequestError(err.Error(), nil)
	}

	return e.JSON(200, map[string]any{
		"success": true,
		"message": fmt.Sprintf("Halbjahresabschluss für %d/%d, %d. Halbjahr erfolgreich erstellt.", year, year+1, semester),
		"year_id": yearRecord.Id,
	})
}

// selectTeachers handles the selection of teachers for a Halbjahresabschluss
func selectTeachers(e *core.RequestEvent) error {
	var req SelectTeachersRequest
	if err := e.BindBody(&req); err != nil {
		return apis.NewBadRequestError("Ungültige Anfrage: "+err.Error(), nil)
	}

	if req.YearID == "" {
		return apis.NewBadRequestError("year_id ist erforderlich", nil)
	}

	if len(req.Teachers) == 0 {
		return apis.NewBadRequestError("Mindestens ein Lehrer muss ausgewählt werden", nil)
	}

	// Validate the year exists and is in the correct state
	year := YearsRecord{}
	err := e.App.DB().
		Select("*").
		From("years").
		Where(dbx.NewExp("id={:id}", dbx.Params{"id": req.YearID})).
		One(&year)
	if err != nil {
		return apis.NewNotFoundError(
			fmt.Sprintf("Halbjahresabschluss mit ID '%s' nicht gefunden", req.YearID),
			map[string]any{"year_id": req.YearID, "error": err.Error()},
		)
	}

	if year.State != YearsRecordStateUploaded {
		return apis.NewBadRequestError(
			fmt.Sprintf("Halbjahresabschluss ist im Status '%s', erwartet wurde 'uploaded'. Lehrerauswahl ist nur für neue Halbjahresabschlüsse möglich.", year.State),
			map[string]any{"current_state": year.State, "expected_state": "uploaded"},
		)
	}

	usersCollection, err := e.App.FindCollectionByNameOrId("users")
	if err != nil {
		return apis.NewApiError(500, "Interne Datenbankfehler: 'users' Collection nicht gefunden", nil)
	}

	tdCollection, err := e.App.FindCollectionByNameOrId("teacher_data")
	if err != nil {
		return apis.NewApiError(500, "Interne Datenbankfehler: 'teacher_data' Collection nicht gefunden", nil)
	}

	// Process each teacher
	createdUsers := []string{}
	createdTD := []string{}
	errors := []string{}

	for name, teacher := range req.Teachers {
		// Validate teacher data
		if len(teacher.Short) < 2 {
			errors = append(errors, fmt.Sprintf("Lehrer '%s': Kürzel '%s' ist zu kurz (min. 2 Zeichen)", name, teacher.Short))
			continue
		}

		if !strings.Contains(teacher.Email, "@") {
			errors = append(errors, fmt.Sprintf("Lehrer '%s': E-Mail '%s' ist ungültig", name, teacher.Email))
			continue
		}

		// Check if user exists by short code
		var userID string
		existingUser := UserRecord{}
		err := e.App.DB().
			Select("id").
			From("users").
			Where(dbx.NewExp("LOWER(short)={:short}", dbx.Params{"short": strings.ToLower(teacher.Short)})).
			One(&existingUser)

		if err == nil && existingUser.Id != "" {
			// User exists, use their ID
			userID = existingUser.Id
		} else {
			// Check if user exists by email
			err = e.App.DB().
				Select("id").
				From("users").
				Where(dbx.NewExp("LOWER(email)={:email}", dbx.Params{"email": strings.ToLower(teacher.Email)})).
				One(&existingUser)

			if err == nil && existingUser.Id != "" {
				userID = existingUser.Id
			} else {
				// Create new user
				userRecord := core.NewRecord(usersCollection)
				userRecord.Set("email", strings.ToLower(teacher.Email))
				userRecord.Set("name", name)
				userRecord.Set("short", strings.ToUpper(teacher.Short))

				if err := e.App.Save(userRecord); err != nil {
					errors = append(errors, fmt.Sprintf("Lehrer '%s': Fehler beim Erstellen des Benutzers: %v", name, err))
					continue
				}
				userID = userRecord.Id
				createdUsers = append(createdUsers, name)
			}
		}

		// Check if teacher_data already exists for this user and semester
		existingTD := TeacherDataRecord{}
		err = e.App.DB().
			Select("id").
			From("teacher_data").
			Where(dbx.NewExp("user={:user} AND semester={:sem}", dbx.Params{"user": userID, "sem": req.YearID})).
			One(&existingTD)

		if err == nil && existingTD.ID != "" {
			errors = append(errors, fmt.Sprintf("Lehrer '%s': Daten für dieses Halbjahr existieren bereits", name))
			continue
		}

		// Create teacher_data record
		tdRecord := core.NewRecord(tdCollection)
		tdRecord.Set("user", userID)
		tdRecord.Set("semester", req.YearID)
		tdRecord.Set("avg_time", teacher.AvgTime)
		tdRecord.Set("class_lead", 0)
		tdRecord.Set("add_points", 0)

		if err := e.App.Save(tdRecord); err != nil {
			errors = append(errors, fmt.Sprintf("Lehrer '%s': Fehler beim Erstellen der Lehrerdaten: %v", name, err))
			continue
		}
		createdTD = append(createdTD, name)
	}

	// If there were any errors, return them but still report success for what worked
	if len(errors) > 0 && len(createdTD) == 0 {
		return apis.NewBadRequestError(
			"Keine Lehrerdaten konnten erstellt werden",
			map[string]any{"errors": errors},
		)
	}

	// Update year state to "open"
	yearsCollection, err := e.App.FindCollectionByNameOrId("years")
	if err != nil {
		return apis.NewApiError(500, "Interne Datenbankfehler: 'years' Collection nicht gefunden", nil)
	}

	yearRecord, err := e.App.FindRecordById(yearsCollection, req.YearID)
	if err != nil {
		return apis.NewApiError(500, "Fehler beim Laden des Halbjahresabschlusses: "+err.Error(), nil)
	}

	yearRecord.Set("state", "open")
	if err := e.App.Save(yearRecord); err != nil {
		return apis.NewApiError(500, "Fehler beim Aktualisieren des Halbjahresabschluss-Status: "+err.Error(), nil)
	}

	response := map[string]any{
		"success":       true,
		"message":       fmt.Sprintf("%d Lehrerdaten erstellt, %d neue Benutzer angelegt", len(createdTD), len(createdUsers)),
		"created_users": createdUsers,
		"created_data":  createdTD,
	}

	if len(errors) > 0 {
		response["warnings"] = errors
	}

	return e.JSON(200, response)
}

// closeYear handles closing a year and triggering the calculation
func closeYear(e *core.RequestEvent) error {
	var req CloseYearRequest
	if err := e.BindBody(&req); err != nil {
		return apis.NewBadRequestError("Ungültige Anfrage: "+err.Error(), nil)
	}

	if req.YearID == "" {
		return apis.NewBadRequestError("year_id ist erforderlich", nil)
	}

	// Validate parameters
	if req.BaseMul <= 0 {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültiger Basisfaktor: %.2f. Der Wert muss größer als 0 sein.", req.BaseMul),
			map[string]any{"field": "base_mul", "value": req.BaseMul},
		)
	}

	if req.LeadPoints < 0 {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültige Klassenleitungspunkte: %.2f. Der Wert darf nicht negativ sein.", req.LeadPoints),
			map[string]any{"field": "lead_points", "value": req.LeadPoints},
		)
	}

	if req.TotalPoints <= 0 {
		return apis.NewBadRequestError(
			fmt.Sprintf("Ungültige Gesamtanrechnungsstunden: %.2f. Der Wert muss größer als 0 sein.", req.TotalPoints),
			map[string]any{"field": "total_points", "value": req.TotalPoints},
		)
	}

	// Validate the year exists and is in the correct state
	year := YearsRecord{}
	err := e.App.DB().
		Select("*").
		From("years").
		Where(dbx.NewExp("id={:id}", dbx.Params{"id": req.YearID})).
		One(&year)
	if err != nil {
		return apis.NewNotFoundError(
			fmt.Sprintf("Halbjahresabschluss mit ID '%s' nicht gefunden", req.YearID),
			map[string]any{"year_id": req.YearID, "error": err.Error()},
		)
	}

	if year.State != YearsRecordStateOpen {
		return apis.NewBadRequestError(
			fmt.Sprintf("Halbjahresabschluss ist im Status '%s', erwartet wurde 'open'. Abschließen ist nur für offene Halbjahresabschlüsse möglich.", year.State),
			map[string]any{"current_state": year.State, "expected_state": "open"},
		)
	}

	// Update teacher_data records
	updatedCount := 0
	for tdID, data := range req.TeacherData {
		if data.ClassLead < 0 || data.ClassLead > 100 {
			return apis.NewBadRequestError(
				fmt.Sprintf("Ungültige Klassenleitung für ID %s: %.2f%%. Der Wert muss zwischen 0 und 100 liegen.", tdID, data.ClassLead),
				map[string]any{"td_id": tdID, "class_lead": data.ClassLead},
			)
		}

		_, err := e.App.DB().
			NewQuery("UPDATE teacher_data SET class_lead={:cl}, add_points={:ap} WHERE id={:id}").
			Bind(dbx.Params{
				"cl": data.ClassLead,
				"ap": data.AddPoints,
				"id": tdID,
			}).
			Execute()
		if err != nil {
			return apis.NewApiError(500,
				fmt.Sprintf("Fehler beim Aktualisieren der Lehrerdaten (ID: %s): %v", tdID, err),
				map[string]any{"td_id": tdID, "error": err.Error()},
			)
		}
		updatedCount++
	}

	// Update year record and trigger calculation
	yearsCollection, err := e.App.FindCollectionByNameOrId("years")
	if err != nil {
		return apis.NewApiError(500, "Interne Datenbankfehler: 'years' Collection nicht gefunden", nil)
	}

	yearRecord, err := e.App.FindRecordById(yearsCollection, req.YearID)
	if err != nil {
		return apis.NewApiError(500, "Fehler beim Laden des Halbjahresabschlusses: "+err.Error(), nil)
	}

	yearRecord.Set("state", "closed")
	yearRecord.Set("base_mul", req.BaseMul)
	yearRecord.Set("lead_points", req.LeadPoints)
	yearRecord.Set("total_points", req.TotalPoints)

	if err := e.App.Save(yearRecord); err != nil {
		return apis.NewApiError(500,
			"Fehler beim Abschließen des Halbjahresabschlusses: "+err.Error(),
			map[string]any{"error": err.Error()},
		)
	}

	// The calculation is triggered by the OnRecordAfterUpdateSuccess hook in years.go
	// But we can also verify it worked by checking for results

	return e.JSON(200, map[string]any{
		"success":       true,
		"message":       fmt.Sprintf("Halbjahresabschluss erfolgreich abgeschlossen. %d Lehrerdaten aktualisiert.", updatedCount),
		"updated_count": updatedCount,
		"calculation_params": map[string]any{
			"base_mul":     req.BaseMul,
			"lead_points":  req.LeadPoints,
			"total_points": req.TotalPoints,
		},
	})
}

// Helper functions

func parseGermanDate(dateStr string) (time.Time, error) {
	// Try various German date formats
	formats := []string{
		"2.1.2006",
		"02.01.2006",
		"2.1.06",
		"02.01.06",
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("konnte Datum '%s' nicht parsen", dateStr)
}

func germanWeekday(w time.Weekday) string {
	switch w {
	case time.Monday:
		return "Montag"
	case time.Tuesday:
		return "Dienstag"
	case time.Wednesday:
		return "Mittwoch"
	case time.Thursday:
		return "Donnerstag"
	case time.Friday:
		return "Freitag"
	case time.Saturday:
		return "Samstag"
	case time.Sunday:
		return "Sonntag"
	}
	return w.String()
}

func hasValidExtension(filename string, allowed []string) bool {
	lower := strings.ToLower(filename)
	for _, ext := range allowed {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
