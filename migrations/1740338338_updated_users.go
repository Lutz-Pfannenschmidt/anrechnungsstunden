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
			"otp": {
				"emailTemplate": {
					"body": "<p>Hallo, </p>\n<p>bitte klicken Sie auf den folgenden Knopf, um sich anzumelden:</p>\n<a class=\"btn\" href=\"{APP_URL}/login/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende Einmalpasswort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn Sie nicht versuchen, sich anzumelden, können Sie diese E-Mail ignorieren.</i></p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
				}
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
			"otp": {
				"emailTemplate": {
					"body": "<p>Hallo, </p>\n<p>bitte klicken Sie auf den folgenden Knopf, um sich anzumelden:</p>\n<a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende Einmalpasswort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn Sie nicht versuchen, sich anzumelden, können Sie diese E-Mail ignorieren.</i></p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
				}
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
