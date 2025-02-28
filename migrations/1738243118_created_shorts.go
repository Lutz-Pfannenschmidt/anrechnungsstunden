package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		jsonData := `{
			"createRule": null,
			"deleteRule": null,
			"fields": [
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "text3208210256",
					"max": 0,
					"min": 0,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"autogeneratePattern": "",
					"hidden": false,
					"id": "_clone_WM7l",
					"max": 0,
					"min": 3,
					"name": "short",
					"pattern": "",
					"presentable": true,
					"primaryKey": false,
					"required": true,
					"system": false,
					"type": "text"
				},
				{
					"exceptDomains": null,
					"hidden": false,
					"id": "_clone_MEDv",
					"name": "email",
					"onlyDomains": null,
					"presentable": false,
					"required": true,
					"system": true,
					"type": "email"
				}
			],
			"id": "pbc_3410128592",
			"indexes": [],
			"listRule": null,
			"name": "shorts",
			"system": false,
			"type": "view",
			"updateRule": null,
			"viewQuery": "SELECT users.short as id, users.short, users.email\nFROM users;",
			"viewRule": null
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3410128592")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
