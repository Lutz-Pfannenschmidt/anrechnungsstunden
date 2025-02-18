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
					"hidden": false,
					"id": "json4224597626",
					"maxSize": 1,
					"name": "subject",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json1499115060",
					"maxSize": 1,
					"name": "grade",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json666537513",
					"maxSize": 1,
					"name": "points",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
				{
					"hidden": false,
					"id": "json2990389176",
					"maxSize": 1,
					"name": "created",
					"presentable": false,
					"required": false,
					"system": false,
					"type": "json"
				},
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
				}
			],
			"id": "pbc_1028955820",
			"indexes": [],
			"listRule": null,
			"name": "current_exam_points",
			"system": false,
			"type": "view",
			"updateRule": null,
			"viewQuery": "WITH LatestPoints AS (\n    SELECT \n        subject,\n        grade,\n        points,\n        created,\n        ROW_NUMBER() OVER (PARTITION BY subject, grade ORDER BY created DESC) AS rn\n    FROM \n        exam_points\n)\nSELECT \n    subject,\n    grade,\n    points,\n    created,\n\t(subject || \"_\" || grade) as id\nFROM \n    LatestPoints\nWHERE \n    rn = 1;\n",
			"viewRule": null
		}`

		collection := &core.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1028955820")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}
