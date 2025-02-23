package main

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func onUserCreate(e *core.RecordEvent) error {
	data := e.Record.PublicExport()

	name := strings.ToLower(data["name"].(string))
	short := strings.ToLower(data["short"].(string))

	e.Record.Set("name", name)
	e.Record.Set("short", short)

	e.App.DB().
		NewQuery("INSERT INTO acronyms (acronym, user) VALUES ({:acro}, {:user}), ({:acro2}, {:user})").
		Bind(dbx.Params{
			"acro":  name,
			"acro2": short,
			"user":  e.Record.Id,
		}).Execute()

	return e.Next()
}

func onUserDelete(e *core.RecordEvent) error {
	data := e.Record.PublicExport()

	name := strings.ToLower(data["name"].(string))
	short := strings.ToLower(data["short"].(string))

	e.App.DB().
		NewQuery("DELETE FROM acronyms WHERE acronym = {:acro} OR acronym = {:acro2}").
		Bind(dbx.Params{
			"acro":  name,
			"acro2": short,
		}).Execute()

	return e.Next()

}

func onUserUpdate(e *core.RecordEvent) error {
	onUserDelete(e)
	onUserCreate(e)

	return e.Next()
}
