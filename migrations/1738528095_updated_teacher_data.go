package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3688353876")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"cascadeDelete": false,
			"collectionId": "_pb_users_auth_",
			"hidden": false,
			"id": "relation2968954581",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "teacher",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "number3145888567",
			"max": null,
			"min": null,
			"name": "year",
			"onlyInt": false,
			"presentable": true,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "number4147678957",
			"max": 2,
			"min": 1,
			"name": "semester",
			"onlyInt": false,
			"presentable": true,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"hidden": false,
			"id": "number2758380978",
			"max": null,
			"min": null,
			"name": "students",
			"onlyInt": false,
			"presentable": true,
			"required": false,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3688353876")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation2968954581")

		// remove field
		collection.Fields.RemoveById("number3145888567")

		// remove field
		collection.Fields.RemoveById("number4147678957")

		// remove field
		collection.Fields.RemoveById("number2758380978")

		return app.Save(collection)
	})
}
