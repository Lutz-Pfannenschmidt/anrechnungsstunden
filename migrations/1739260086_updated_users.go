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
					"body": "<p>Hallo,</p>\n<p>bitte klicken auf den folgenden Knopf, um sich anzumelden:</p>\n  <a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende einmal Passwort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn sie nicht versuchen, sich anzumelden können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
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
					"body": "<p>Hello,</p>\n<p>Your one-time password is: <strong>{OTP}</strong></p>\n  <a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Login</a>\n<p><i>If you didn't ask for the one-time password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>"
				}
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
