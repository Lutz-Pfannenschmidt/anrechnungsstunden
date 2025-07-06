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

func onUserCreateBefore(e *core.RecordEvent) error {
	e.Record.Set("short", strings.ToLower(e.Record.GetString("short")))
	e.Record.Set("email", strings.ToLower(e.Record.GetString("email")))
	return e.Next()
}

func pdfSender(e *core.RequestEvent) error {
	year_id := e.Request.URL.Query().Get("year")
	year := YearsRecord{}
	err := e.App.DB().
		Select("*").
		From("years").
		AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": year_id})).
		One(&year)
	if err != nil {
		return fmt.Errorf("failed to retrieve year: %w", err)
	}

	if year.ID != year_id {
		return e.Error(404, "Year not found", nil)
	}

	pdf_rows := []PdfsRecord{}

	err = e.App.DB().
		Select("*").
		From("pdfs").
		AndWhere(dbx.NewExp("semester={:semester}", dbx.Params{"semester": year_id})).
		All(&pdf_rows)
	if err != nil {
		return fmt.Errorf("failed to retrieve PDF records: %w", err)
	}

	registry := template.NewRegistry()
	for _, row := range pdf_rows {
		u := UserRecord{}
		err = e.App.DB().
			Select("*").
			From("users").
			AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": row.User})).
			One(&u)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %w", err)
		}

		yearStr := fmt.Sprintf("%d. Halbjahr %d/%s", year.Semester, year.StartYear, strconv.Itoa(year.StartYear + 1)[2:])
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

	return nil
}
