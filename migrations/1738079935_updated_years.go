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

		// remove field
		collection.Fields.RemoveById("bool17505815")

		// remove field
		collection.Fields.RemoveById("bool3680925218")

		// add field
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
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1749964444")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "bool17505815",
			"name": "unlocked",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"hidden": false,
			"id": "bool3680925218",
			"name": "parsed",
			"presentable": true,
			"required": false,
			"system": false,
			"type": "bool"
		}`)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("select2744374011")

		return app.Save(collection)
	})
}
