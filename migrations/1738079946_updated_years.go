package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1749964444")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"hidden": false,
			"id": "select2744374011",
			"maxSelect": 1,
			"name": "state",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "select",
			"values": [
				"uploaded",
				"parsed",
				"open",
				"closed"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1749964444")
		if err != nil {
			return err
		}

		// update field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"hidden": false,
			"id": "select2744374011",
			"maxSelect": 1,
			"name": "state",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "select",
			"values": [
				"uploaded",
				"parsed",
				"open",
				"closed"
			]
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
