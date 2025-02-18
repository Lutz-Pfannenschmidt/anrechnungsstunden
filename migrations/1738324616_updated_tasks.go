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
			"listRule": "  @request.auth.id = user.id",
			"viewQuery": "SELECT (ROW_NUMBER() OVER()) as id, years.id AS year, users.id AS user\nFROM years\nJOIN users ON years.must_complete LIKE '%\"' || users.id || '\"%' AND years.state = \"open\";\n",
			"viewRule": "  @request.auth.id = user.id"
		}`), &collection); err != nil {
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
			"listRule": null,
			"viewQuery": "SELECT (ROW_NUMBER() OVER()) as id, years.id AS year, users.id AS user\nFROM years\nJOIN users ON years.must_complete LIKE '%\"' || users.id || '\"%';",
			"viewRule": null
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
