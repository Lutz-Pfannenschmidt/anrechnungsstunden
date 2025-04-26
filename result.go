package main

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func onResultsDelete(e *core.RecordEvent) error {
	semID := e.Record.GetString("semester")
	if semID == "" {
		return nil
	}

	sem, err := e.App.FindRecordById("years", semID)
	if err != nil {
		return err
	}

	sem.Set("state", "open")
	err = e.App.Save(sem)
	if err != nil {
		return err
	}

	pdfs, err := e.App.FindAllRecords("pdfs", dbx.HashExp{"semester": semID})
	if err != nil {
		return err
	}

	for _, pdf := range pdfs {
		err = e.App.Delete(pdf)
		if err != nil {
			return err
		}
	}

	return nil
}
