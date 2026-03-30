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
		return fmt.Errorf("Halbjahresabschluss '%s' konnte nicht abgerufen werden: %w", year_id, err)
	}

	if year.ID != year_id {
		return e.Error(404, "Halbjahresabschluss nicht gefunden", nil)
	}

	pdf_rows := []PdfsRecord{}

	err = e.App.DB().
		Select("*").
		From("pdfs").
		AndWhere(dbx.NewExp("semester={:semester}", dbx.Params{"semester": year_id})).
		All(&pdf_rows)
	if err != nil {
		return fmt.Errorf("PDF-Datensätze konnten nicht abgerufen werden: %w", err)
	}

	registry := template.NewRegistry()
	for _, row := range pdf_rows {
		if err := sendPdfEmail(e, registry, &year, &row); err != nil {
			return err
		}
	}

	return nil
}

// sendPdfEmail sends a single PDF email to a user.
// Extracted from pdfSender to properly handle file closing.
func sendPdfEmail(e *core.RequestEvent, registry *template.Registry, year *YearsRecord, row *PdfsRecord) error {
	u := UserRecord{}
	err := e.App.DB().
		Select("*").
		From("users").
		AndWhere(dbx.NewExp("id={:id}", dbx.Params{"id": row.User})).
		One(&u)
	if err != nil {
		return fmt.Errorf("Benutzer '%s' konnte nicht abgerufen werden: %w", row.User, err)
	}

	yearStr := fmt.Sprintf("%d. Halbjahr %d/%s", year.Semester, year.StartYear, strconv.Itoa(year.StartYear + 1)[2:])
	data := map[string]any{
		"yearStr": yearStr,
	}

	html, err := registry.LoadFiles("templates/pdf_mail.html").Render(data)
	if err != nil {
		return fmt.Errorf("E-Mail-Vorlage konnte nicht gerendert werden: %w", err)
	}

	pdfRecord, err := e.App.FindRecordById("pdfs", row.Id)
	if err != nil {
		return fmt.Errorf("PDF-Datensatz '%s' nicht gefunden: %w", row.Id, err)
	}

	pdfPath := e.App.DataDir() + "/storage/" + pdfRecord.BaseFilesPath() + "/" + row.Pdf
	pdf, err := os.Open(pdfPath)
	if err != nil {
		return fmt.Errorf("PDF-Datei '%s' konnte nicht geöffnet werden: %w", pdfPath, err)
	}
	defer pdf.Close()

	message := &mailer.Message{
		From: mail.Address{
			Address: e.App.Settings().Meta.SenderAddress,
			Name:    e.App.Settings().Meta.SenderName,
		},
		To:          []mail.Address{{Address: u.Email, Name: u.Name}},
		Subject:     "Ihre Anrechnungsstunden",
		HTML:        html,
		Attachments: map[string]io.Reader{strings.ToLower(u.Name) + ".pdf": pdf},
	}

	err = e.App.NewMailClient().Send(message)
	if err != nil {
		return fmt.Errorf("E-Mail an '%s' (%s) konnte nicht gesendet werden: %w", u.Name, u.Email, err)
	}

	return nil
}
