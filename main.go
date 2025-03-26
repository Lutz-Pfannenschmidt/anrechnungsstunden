package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "anrechnungsstundenberechner/migrations"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	app.OnRecordAfterCreateSuccess("users").BindFunc(onUserCreate)
	app.OnRecordAfterDeleteSuccess("users").BindFunc(onUserDelete)
	app.OnRecordAfterUpdateSuccess("users").BindFunc(onUserUpdate)

	app.OnRecordAfterCreateSuccess("results").BindFunc(makePdf)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/parse/", parse).BindFunc(apis.RequireSuperuserAuth().Func)
		se.Router.GET("/send_pdfs/", pdfSender).BindFunc(apis.RequireSuperuserAuth().Func)

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
}
