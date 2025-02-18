package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_4089960912")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewQuery": "SELECT \n    (ROW_NUMBER() OVER()) AS id, \n    years.start_year AS year, \n    years.semester AS semester, \n    users.id AS user, \n    time_data.avg_time\nFROM \n    years\nJOIN \n    users ON years.must_complete LIKE '%\"' || users.id || '\"%' \n    AND years.state = 'open'\nJOIN \n    time_data ON time_data.user = users.id;"
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("relation3145888567")

		// remove field
		collection.Fields.RemoveById("_clone_hpxl")

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "_clone_LCVB",
			"max": null,
			"min": 2000,
			"name": "year",
			"onlyInt": true,
			"presentable": true,
			"required": true,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(2, []byte(`{
			"hidden": false,
			"id": "_clone_eqCa",
			"max": 2,
			"min": 1,
			"name": "semester",
			"onlyInt": true,
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
			"id": "_clone_1Tj1",
			"max": null,
			"min": null,
			"name": "avg_time",
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
		collection, err := app.FindCollectionByNameOrId("pbc_4089960912")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewQuery": "SELECT \n    (ROW_NUMBER() OVER()) AS id, \n    years.id AS year, \n    users.id AS user, \n    time_data.avg_time\nFROM \n    years\nJOIN \n    users ON years.must_complete LIKE '%\"' || users.id || '\"%' \n    AND years.state = 'open'\nJOIN \n    time_data ON time_data.user = users.id;"
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"cascadeDelete": false,
			"collectionId": "pbc_1749964444",
			"hidden": false,
			"id": "relation3145888567",
			"maxSelect": 1,
			"minSelect": 0,
			"name": "year",
			"presentable": false,
			"required": false,
			"system": false,
			"type": "relation"
		}`)); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(3, []byte(`{
			"hidden": false,
			"id": "_clone_hpxl",
			"max": null,
			"min": null,
			"name": "avg_time",
			"onlyInt": false,
			"presentable": true,
			"required": false,
			"system": false,
			"type": "number"
		}`)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("_clone_LCVB")

		// remove field
		collection.Fields.RemoveById("_clone_eqCa")

		// remove field
		collection.Fields.RemoveById("_clone_1Tj1")

		return app.Save(collection)
	})
}
