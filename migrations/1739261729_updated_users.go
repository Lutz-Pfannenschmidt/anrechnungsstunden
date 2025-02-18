package migrations

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"resetPasswordTemplate": {
				"body": "<p>Hallo,</p>\n<p>zum zurücksetzen ihres Passwort klicken sie btte auf den folgenden Knopf.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Passwort zurücksetzen</a>\n</p>\n<p><i>Wenn sie nicht versucht haben, ihr Passwort zurückzusetzen, können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// update collection data
		if err := json.Unmarshal([]byte(`{
			"resetPasswordTemplate": {
				"body": "<p>Hallo,</p>\n<p>zum zurücksetzen ihres Passwort klicken sie btte auf den folgenden Knopf.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>Wenn sie nicht versucht haben, ihr Passwort zurückzusetzen, können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
