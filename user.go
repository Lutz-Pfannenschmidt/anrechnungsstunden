package main

import (
	"fmt"
	"io"
	"net/mail"
	"os"
	"strconv"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/mailer"
	"github.com/pocketbase/pocketbase/tools/template"
)

type Pdf_Row struct {
	Id       string `db:"id" json:"id"`
	Semester string `db:"semester" json:"semester"`
	User     string `db:"user" json:"user"`
	Pdf      string `db:"pdf" json:"pdf"`
}

func onUserCreate(e *core.RecordEvent) error {
	data := e.Record.PublicExport()

	short := strings.ToLower(data["short"].(string))

	e.App.DB().
		NewQuery("INSERT INTO acronyms (acronym, user) VALUES ({:acro}, {:user})").
		Bind(dbx.Params{
			"acro": short,
			"user": e.Record.Id,
		}).Execute()

	return e.Next()
}

func onUserDelete(e *core.RecordEvent) error {
	data := e.Record.PublicExport()

	short := strings.ToLower(data["short"].(string))

	e.App.DB().
		NewQuery("DELETE FROM acronyms WHERE acronym = {:acro}").
		Bind(dbx.Params{
			"acro": short,
		}).Execute()

	return e.Next()

}

func onUserUpdate(e *core.RecordEvent) error {
	onUserDelete(e)
	onUserCreate(e)

	return e.Next()
}

func pdfSender(e *core.RequestEvent) error {
	year_id := e.Request.URL.Query().Get("year")
	year := Year{}
	e.App.DB().
		Select("*").
		From("years").
		AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": year_id})).
		One(&year)

	if year.ID != year_id {
		return e.Error(404, "Year not found", nil)
	}

	pdf_rows := []Pdf_Row{}

	e.App.DB().
		Select("*").
		From("pdfs").
		AndWhere(dbx.NewExp("semester={:semester}", dbx.Params{"semester": year_id})).
		All(&pdf_rows)

	registry := template.NewRegistry()
	for _, row := range pdf_rows {
		u := User{}
		e.App.DB().
			Select("*").
			From("users").
			AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": row.User})).
			One(&u)

		yearStr := fmt.Sprintf("%d. Halbjahr %d/%s", year.Semester, year.Year, strconv.Itoa(year.Year + 1)[2:])
		data := map[string]any{
			"yearStr": yearStr,
		}

		html, err := registry.LoadFiles("templates/pdf_mail.html").Render(data)
		if err != nil {
			return e.Error(500, "Error rendering template", err)
		}

		pdfRecord, _ := e.App.FindRecordById("pdfs", row.Id)
		pdf, err := os.Open(e.App.DataDir() + "/storage/" + pdfRecord.BaseFilesPath() + "/" + row.Pdf)
		if err != nil {
			return e.Error(500, "Error opening pdf", err)
		}
		defer pdf.Close()

		message := &mailer.Message{
			From: mail.Address{
				Address: e.App.Settings().Meta.SenderAddress,
				Name:    e.App.Settings().Meta.SenderName,
			},
			To:          []mail.Address{{Address: u.Email, Name: u.Name}},
			Subject:     "Ihre Anrechnugsstunden",
			HTML:        html,
			Attachments: map[string]io.Reader{strings.ToLower(u.Name) + ".pdf": pdf},
		}

		err = e.App.NewMailClient().Send(message)
		if err != nil {
			return e.Error(500, "Error sending mail", err)
		}
	}

	return e.Next()
}
