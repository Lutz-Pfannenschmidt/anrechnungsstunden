package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3172483855")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_GoZX9LAQqa` + "`" + ` ON ` + "`" + `exam_points` + "`" + ` (\n  ` + "`" + `grade` + "`" + `,\n  ` + "`" + `lesson` + "`" + `\n)"
			],
			"name": "exam_points"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3172483855")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_GoZX9LAQqa` + "`" + ` ON ` + "`" + `lesson_points` + "`" + ` (\n  ` + "`" + `grade` + "`" + `,\n  ` + "`" + `lesson` + "`" + `\n)"
			],
			"name": "lesson_points"
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
