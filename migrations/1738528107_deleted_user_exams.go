package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1967380704")
		if err != nil {
			return nil
		}

		return app.Delete(collection)
	}, func(app core.App) error {
		jsonData := `{
			"createRule": "@request.auth.id != \"\"",
			"deleteRule": "@request.auth.id = teacher.id",
			"fields": [
				{
					"autogeneratePattern": "[a-z0-9]{15}",
					"hidden": false,
					"id": "text3208210256",
					"max": 15,
					"min": 15,
					"name": "id",
					"pattern": "^[a-z0-9]+$",
					"presentable": false,
					"primaryKey": true,
					"required": true,
					"system": true,
					"type": "text"
				},
				{
					"cascadeDelete": false,
					"collectionId": "pbc_3172483855",
					"hidden": false,
					"id": "relation3981121951",
					"maxSelect": 1,
					"minSelect": 0,
					"name": "class",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "relation"
				},
				{
					"hidden": false,
					"id": "number2758380978",
					"max": null,
					"min": null,
					"name": "students",
					"onlyInt": true,
					"presentable": true,
					"required": true,
					"system": false,
					"type": "number"
				},
				{
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
				},
				{
					"hidden": false,
					"id": "autodate2990389176",
					"name": "created",
					"onCreate": true,
					"onUpdate": false,
					"presentable": false,
					"system": false,
					"type": "autodate"
				},
				{
					"hidden": false,
					"id": "autodate3332085495",
					"name": "updated",
					"onCreate": true,
					"onUpdate": true,
					"presentable": false,
					"system": false,
					"type": "autodate"
				}
			],
			"id": "pbc_1967380704",
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_J5SSL33PnD` + "`" + ` ON ` + "`" + `user_exams` + "`" + ` (\n  ` + "`" + `teacher` + "`" + `,\n  ` + "`" + `class` + "`" + `\n)"
			],
			"listRule": "@request.auth.id = teacher.id",
			"name": "user_exams",
			"system": false,
			"type": "base",
			"updateRule": "@request.auth.id = teacher.id",
			"viewRule": "@request.auth.id = teacher.id"
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
