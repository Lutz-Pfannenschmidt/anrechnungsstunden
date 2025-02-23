package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"anrechnungsstundenberechner/internal/out"
	"anrechnungsstundenberechner/internal/parser"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Year struct {
	ID           string                  `db:"id" json:"id"`
	File         string                  `db:"file" json:"file"`
	Semester     int                     `db:"semester" json:"semester"`
	Year         int                     `db:"start_year" json:"start_year"`
	MustComplete types.JSONArray[string] `db:"must_complete" json:"must_complete"`
}

type Result struct {
	Id         string                 `db:"id" json:"id"`
	Semester   string                 `db:"semester" json:"semester"`
	Data       types.JSONMap[float64] `db:"data" json:"data"`
	Pdf        string                 `db:"pdf" json:"pdf"`
	LeadPoints float64                `db:"lead_points" json:"lead_points"`
}

func (r Result) toSring() string {
	return fmt.Sprintf("Id: %s, Semester: %s, Data: %v, Pdf: %s, LeadPoints: %f", r.Id, r.Semester, r.Data, r.Pdf, r.LeadPoints)
}

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Verified bool   `db:"verified" json:"verified"`
	Name     string `db:"name" json:"name"`
	Avatar   string `db:"avatar" json:"avatar"`
	Short    string `db:"short" json:"short"`
}

type TeacherData struct {
	Id       string `db:"id" json:"id"`
	Teacher  string `db:"teacher" json:"teacher"`
	Year     int    `db:"year" json:"year"`
	Semester int    `db:"semester" json:"semester"`
	Students int    `db:"students" json:"students"`
	Grade    string `db:"grade" json:"grade"`
	Subject  string `db:"subject" json:"subject"`
}

type currentExamPoints struct {
	Id      string `db:"id" json:"id"`
	Subject string `db:"subject" json:"subject"`
	Grade   string `db:"grade" json:"grade"`
	Points  int    `db:"points" json:"points"`
}

const templatePath = "templates/template.xlsx"

var cleanup []string = []string{}

func makePdf(e *core.RecordEvent) error {
	var r Result
	data, _ := json.Marshal(e.Record.PublicExport())
	json.Unmarshal(data, &r)

	fmt.Printf("result: %s\n", r.toSring())

	y := Year{}
	e.App.DB().
		Select("*").
		From("years").
		AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": r.Semester})).
		One(&y)

	fmt.Printf("year: %d, semester: %d", y.Year, y.Semester)
	fmt.Printf("must complete: %v", y.MustComplete)

	yearStr := fmt.Sprintf("%d/%s, %d. Halbjahr", y.Year, strconv.Itoa(y.Year + 1)[2:], y.Semester)

	files := []string{}
	for _, uid := range y.MustComplete {
		fmt.Printf("uid: %s\n", uid)
		type Res struct {
			Id                  string  `db:"id"`
			User                string  `db:"user"`
			Name                string  `db:"name"`
			Short               string  `db:"short"`
			AvgTime             float64 `db:"avg_time"`
			ClassLeadPercentage float64 `db:"class_lead_percentage"`
		}

		user_data := Res{}
		e.App.DB().NewQuery(
			`SELECT 
                    (ROW_NUMBER() OVER()) AS id,
                    users.id AS user, 
                    users.name,
                    users.short,
                    time_data.avg_time,
                    class_lead.percentage as class_lead_percentage
                FROM 
                    years
                JOIN 
                    users ON years.must_complete LIKE '%"' || {:uid} || '"%' 
                    AND years.state = 'open'
                JOIN 
                    time_data ON time_data.user = users.id AND time_data.semester = years.id
                JOIN
                    class_lead ON years.start_year = class_lead.year AND years.semester = class_lead.semester AND class_lead.teacher = users.id
                WHERE
                    years.id = {:yearId} AND users.id = {:uid}`).
			Bind(dbx.Params{"uid": uid, "yearId": y.ID}).
			One(&user_data)

		if user_data.User == "" {
			continue
		}

		teacher_data := []TeacherData{}
		e.App.DB().NewQuery(`
                SELECT * FROM teacher_data
                WHERE teacher_data.year={:year} AND teacher_data.semester={:semester}
                AND teacher_data.teacher={:uid}`).
			Bind(dbx.Params{"uid": uid, "year": y.Year, "semester": y.Semester}).
			All(&teacher_data)

		data := []out.RowData{}

		for _, td := range teacher_data {
			p := currentExamPoints{}
			err := e.App.DB().NewQuery(`
					SELECT points FROM current_exam_points
					WHERE subject={:subject} AND grade={:grade}`).
				Bind(dbx.Params{"subject": td.Subject, "grade": td.Grade}).
				One(&p)
			if err != nil {
				p.Points = -1
				fmt.Printf("err: %v\n", err)
			}

			data = append(data, out.RowData{
				Subject:  td.Subject,
				Grade:    td.Grade,
				Students: td.Students,
				Points:   p.Points,
			})
		}

		path := os.TempDir() + "/" + out.TempName(fmt.Sprintf("-%s-%s.pdf", user_data.User, user_data.Id))
		files = append(files, path)
		out.RenderTemplate(
			templatePath,
			path,
			yearStr,
			user_data.Name,
			user_data.AvgTime,
			int(user_data.ClassLeadPercentage),
			r.LeadPoints,
			r.Data[user_data.Name],
			data,
		)

		collection, err := e.App.FindCollectionByNameOrId("pdfs")
		if err != nil {
			return err
		}

		f, err := filesystem.NewFileFromPath(path)
		if err != nil {
			return fmt.Errorf("error creating file from path: %w", err)
		}

		record := core.NewRecord(collection)
		record.Set("user", user_data.User)
		record.Set("semester", y.ID)
		record.Set("pdf", f)

		err = e.App.Save(record)
		if err != nil {
			return err
		}
	}

	dirname := os.TempDir() + "/" + out.TempName("/")
	os.MkdirAll(dirname, os.ModePerm)
	outPath := dirname + "ausgabe.pdf"
	mergePdf(files, outPath)

	f, err := filesystem.NewFileFromPath(outPath)
	if err != nil {
		return fmt.Errorf("error creating file from path: %w", err)
	}
	e.Record.Set("pdf", f)

	for _, path := range files {
		os.Remove(path)
	}

	cleanup = append(cleanup, outPath, dirname)

	return nil
}

func mergePdf(inputPaths []string, outputPath string) error {
	config := model.NewDefaultConfiguration()
	err := api.MergeCreateFile(inputPaths, outputPath, false, config)
	if err != nil {
		return fmt.Errorf("error merging PDFs: %w", err)
	}
	return nil
}

func parse(e *core.RequestEvent) error {
	year := e.Request.URL.Query().Get("year")
	semester := e.Request.URL.Query().Get("semester")

	semInt, err := strconv.Atoi(semester)
	if err != nil {
		return err
	}

	y := Year{}

	e.App.DB().
		Select("*").
		From("years").
		AndWhere(dbx.NewExp("start_year = {:year} AND semester = {:semester}", dbx.Params{"year": year, "semester": semester})).
		One(&y)

	record, err := e.App.FindRecordById("years", y.ID)
	if err != nil {
		return err
	}

	path := e.App.DataDir() + "/storage/" + record.BaseFilesPath() + "/" + record.GetString("file")

	data, err := parser.ParseFile(path)
	if err != nil {
		return err
	}

	filtered := map[string]float64{}
	for k, v := range data {
		filtered[k] = v[semInt-1]
	}

	e.JSON(200, filtered)
	return nil
}
