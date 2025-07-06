package main

import (
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func onYearsDelete(e *core.RecordEvent) error {
	yearID := e.Record.Id
	year := e.Record.GetInt("start_year")
	semester := e.Record.GetInt("semester")
	if yearID == "" || year == 0 || semester == 0 {
		return fmt.Errorf("invalid year record: %s with year: %v and sem: %v", e.Record.Id, year, semester)
	}

	// Delete all files associated with the year
	_, err := e.App.DB().
		NewQuery("DELETE FROM files WHERE year = {:year} AND semester = {:sem}").
		Bind(dbx.Params{
			"year": year,
			"sem":  semester,
		}).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete files for year %s: %w", yearID, err)
	}

	return e.Next()
}

func onYearsUpdate(e *core.RecordEvent) error {
	yearID := e.Record.Id

	if e.Record.GetString("state") != "closed" {
		// If the year is not closed, there shouldnt be any saved results or PDFs
		_, err := e.App.DB().
			NewQuery("DELETE FROM results WHERE semester = {:yearID}").
			Bind(dbx.Params{"yearID": yearID}).
			Execute()
		if err != nil {
			return fmt.Errorf("failed to delete results for year %s: %w", yearID, err)
		}

		_, err = e.App.DB().
			NewQuery("DELETE FROM pdfs WHERE semester = {:yearID}").
			Bind(dbx.Params{"yearID": yearID}).
			Execute()
		if err != nil {
			return fmt.Errorf("failed to delete PDFs for year %s: %w", yearID, err)
		}
	}

	resRecord := ResultRecord{}
	err := e.App.DB().
		NewQuery("SELECT * FROM results WHERE semester = {:yearID}").
		Bind(dbx.Params{"yearID": yearID}).
		One(&resRecord)
	if err != nil {
		resRecord.Id = ""
	}

	if e.Record.GetString("state") == "open" && resRecord.Id != "" {
		// If the year is opened, we delete the result record if it exists
		_, err = e.App.DB().
			NewQuery("DELETE FROM results WHERE id = {:id}").
			Bind(dbx.Params{"id": resRecord.Id}).
			Execute()
		if err != nil {
			return fmt.Errorf("failed to delete result record for year %s: %w", yearID, err)
		}
		resRecord.Id = ""
	}

	if e.Record.GetString("state") == "closed" && resRecord.Id == "" {
		return CalculateAndRenderResults(e.App, yearID)
	}

	return nil
}
