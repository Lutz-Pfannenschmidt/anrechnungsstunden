package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1028955820")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewQuery": "WITH LatestPoints AS (\n    SELECT \n        subject,\n        grade,\n        points,\n        created,\n        ROW_NUMBER() OVER (PARTITION BY subject, grade ORDER BY created DESC) AS rn,\n        id\n    FROM \n        exam_points\n)\nSELECT \n\tid,\n    subject,\n    grade,\n    points,\n    created\nFROM \n    LatestPoints\nWHERE \n    rn = 1;\n"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_1028955820")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"viewQuery": "WITH LatestPoints AS (\n    SELECT \n        subject,\n        grade,\n        points,\n        created,\n        ROW_NUMBER() OVER (PARTITION BY subject, grade ORDER BY created DESC) AS rn,\n        id as old_id\n    FROM \n        exam_points\n)\nSELECT \n\told_id as id,\n    subject,\n    grade,\n    points,\n    created\nFROM \n    LatestPoints\nWHERE \n    rn = 1;\n"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
