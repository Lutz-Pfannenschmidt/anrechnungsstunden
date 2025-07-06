package main

import (
	"fmt"
	"net/http"

	"anrechnungsstundenberechner/internal/parser"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func parse(e *core.RequestEvent) error {
	reqData := struct {
		Year      uint   `json:"year"`       // YYYY
		Semester  uint8  `json:"semester"`   // 1 or 2
		SplitDate string `json:"split_date"` // DD.MM.YYYY
	}{}
	if err := e.BindBody(&reqData); err != nil {
		return e.BadRequestError("Failed to read request data", err)
	}

	hours_file_record, err := e.App.FindFirstRecordByFilter("files", "year = {:year} && semester = {:semester} && type = 'hours'", dbx.Params{"year": reqData.Year, "semester": reqData.Semester})
	if err != nil {
		return err
	}

	path := e.App.DataDir() + "/storage/" + hours_file_record.BaseFilesPath() + "/" + hours_file_record.GetString("file")
	data, err := parser.ParseFile(path, int(reqData.Year), reqData.SplitDate)
	if err != nil {
		return e.Error(http.StatusInternalServerError, "Failed to parse file at "+path, err)
	}

	filtered := map[string]float64{}
	for k, v := range data.Result {
		filtered[k] = v[reqData.Semester-1]
	}

	err = e.JSON(200, map[string]any{
		"result":          filtered,
		"name_collisions": data.NameCollisions,
	})
	if err != nil {
		return fmt.Errorf("failed to send JSON response: %w", err)
	}
	return nil
}
