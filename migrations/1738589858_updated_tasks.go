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
			"viewQuery": "SELECT \n    (ROW_NUMBER() OVER()) AS id, \n    years.start_year AS year, \n    years.semester AS semester, \n    users.id AS user, \n    time_data.avg_time,\n    time_data.created as created\nFROM \n    years\nJOIN \n    users ON years.must_complete LIKE '%\"' || users.id || '\"%' \n    AND years.state = 'open'\nJOIN \n    time_data ON time_data.user = users.id AND time_data.semester = years.id"
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("_clone_EOM2")

		// remove field
		collection.Fields.RemoveById("_clone_PA4Q")

		// remove field
		collection.Fields.RemoveById("_clone_3WOQ")

		// remove field
		collection.Fields.RemoveById("_clone_1hVk")

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "_clone_GhmW",
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
			"id": "_clone_Jvrj",
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
			"id": "_clone_y8e3",
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

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"hidden": false,
			"id": "_clone_itpf",
			"name": "created",
			"onCreate": true,
			"onUpdate": false,
			"presentable": false,
			"system": false,
			"type": "autodate"
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
			"viewQuery": "SELECT \n    (ROW_NUMBER() OVER()) AS id, \n    years.start_year AS year, \n    years.semester AS semester, \n    users.id AS user, \n    time_data.avg_time,\n    time_data.created as created\nFROM \n    years\nJOIN \n    users ON years.must_complete LIKE '%\"' || users.id || '\"%' \n    AND years.state = 'open'\nJOIN \n    time_data ON time_data.user = users.id;"
		}`), &collection); err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(1, []byte(`{
			"hidden": false,
			"id": "_clone_EOM2",
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
			"id": "_clone_PA4Q",
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
			"id": "_clone_3WOQ",
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

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(5, []byte(`{
			"hidden": false,
			"id": "_clone_1hVk",
			"name": "created",
			"onCreate": true,
			"onUpdate": false,
			"presentable": false,
			"system": false,
			"type": "autodate"
		}`)); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("_clone_GhmW")

		// remove field
		collection.Fields.RemoveById("_clone_Jvrj")

		// remove field
		collection.Fields.RemoveById("_clone_y8e3")

		// remove field
		collection.Fields.RemoveById("_clone_itpf")

		return app.Save(collection)
	})
}
