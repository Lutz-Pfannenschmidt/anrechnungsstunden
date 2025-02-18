package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/Lutz-Pfannenschmidt/stunden-berechner/migrations"

	"github.com/Lutz-Pfannenschmidt/stunden-berechner/internal/parser"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func makeParser(app *pocketbase.PocketBase) func(e *core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		year := e.Request.URL.Query().Get("year")
		semester := e.Request.URL.Query().Get("semester")

		semInt, err := strconv.Atoi(semester)
		if err != nil {
			return err
		}

		y := Year{}

		app.DB().
			Select("*").
			From("years").
			AndWhere(dbx.NewExp("start_year = {:year} AND semester = {:semester}", dbx.Params{"year": year, "semester": semester})).
			One(&y)

		record, err := app.FindRecordById("years", y.ID)
		if err != nil {
			return err
		}

		path := app.DataDir() + "/storage/" + record.BaseFilesPath() + "/" + record.GetString("file")

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
}

func main() {
	app := pocketbase.New()

	app.OnRecordCreate("users").BindFunc(func(e *core.RecordEvent) error {
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
	})

	app.OnRecordValidate("results").BindFunc(makePdf)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		se.Router.GET("/parse/", makeParser(app)).BindFunc(apis.RequireSuperuserAuth().Func)

		se.Router.GET("/{path...}", apis.Static(os.DirFS("./client/dist/"), false))
		return se.Next()
	})

	migrate := strings.HasPrefix(filepath.Dir(os.Args[0]), "/tmp") || strings.HasPrefix(filepath.Dir(os.Args[0]), os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: migrate,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}

	for _, fname := range cleanup {
		_ = os.Remove(fname)
	}
}
