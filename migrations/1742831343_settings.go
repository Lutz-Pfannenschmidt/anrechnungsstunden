package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		settings := app.Settings()

		settings.Batch.Enabled = true
		settings.Batch.MaxRequests = 1000
		settings.Batch.MaxBodySize = 512000000
		settings.Batch.Timeout = 30

		return app.Save(settings)
	}, nil)
}
