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
			"viewQuery": "SELECT \n    (ROW_NUMBER() OVER()) AS id, \n    years.id AS year, \n    users.id AS user, \n    time_data.avg_time\nFROM \n    years\nJOIN \n    users ON years.must_complete LIKE '%\"' || users.id || '\"%' \n    AND years.state = 'open'\nJOIN \n    time_data ON time_data.user = users.id;"
		}`), &collection); err != nil {
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

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_4089960912")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewQuery": "SELECT (ROW_NUMBER() OVER()) as id, years.id AS year, users.id AS user\nFROM years\nJOIN users ON years.must_complete LIKE '%\"' || users.id || '\"%' AND years.state = \"open\";\n"
		}`), &collection); err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("_clone_hpxl")

		return app.Save(collection)
	})
}
