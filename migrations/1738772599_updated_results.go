package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_250649598")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "file250665868",
			"maxSelect": 1,
			"maxSize": 0,
			"mimeTypes": [
				"application/pdf"
			],
			"name": "pdf",
			"presentable": false,
			"protected": true,
			"required": false,
			"system": false,
			"thumbs": [],
			"type": "file"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_250649598")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("file250665868")

		return app.Save(collection)
	})
}
