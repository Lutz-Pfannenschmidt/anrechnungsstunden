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
			"authAlert": {
				"emailTemplate": {
					"body": "<p>Hallo, </p>\n<p>Wir haben festgestellt, dass Sie sich von einem neuen Standort aus angemeldet haben</p>\n<p><strong>Wenn Sie das nicht waren, sollten Sie sofort das Passwort für Ihr Konto ändern, um den Zugriff von allen anderen Standorten zu widerrufen. Andernfalls können sie diese E-Mail ignorieren.</strong></p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
				}
			},
			"confirmEmailChangeTemplate": {
				"body": "<p>Hallo, </p>\n<p>klicken Sie bitte auf den folgenden Knopf, um Ihre neue E-Mail-Adresse zu bestätigen.</p>\n<p>\n    <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">E-Mail Bestätigen</a>\n</p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			},
			"otp": {
				"emailTemplate": {
					"body": "<p>Hallo, </p>\n<p>bitte klicken Sie auf den folgenden Knopf, um sich anzumelden:</p>\n<a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende Einmalpasswort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn Sie nicht versuchen, sich anzumelden, können Sie diese E-Mail ignorieren.</i></p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
				}
			},
			"resetPasswordTemplate": {
				"body": "<p>Hallo, </p>\n<p>zum Zurücksetzen ihres Passwortes klicken Sie bitte auf den folgenden Knopf.</p>\n<p>\n    <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Passwort zurücksetzen</a>\n</p>\n<p><i>Wenn sie nicht versucht haben, ihr Passwort zurückzusetzen, können sie diese E-Mail ignorieren.</i></p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			},
			"verificationTemplate": {
				"body": "<p>Hallo, </p>\n<p>bitte klicken Sie auf den folgenden Knopf, um ihre E-Mail-Adresse zu bestätigen.</p>\n<p>\n    <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Bestätigen</a>\n</p>\n<p>\n    Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden Sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
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
			"authAlert": {
				"emailTemplate": {
					"body": "<p>Hallo,</p>\n<p>Wir haben festgestellt, dass Sie sich von einem neuen Standort aus angemeldet haben</p>\n<p><strong>Wenn Sie das nicht waren, sollten Sie sofort das Passwort für Ihr Konto ändern, um den Zugriff von allen anderen Standorten zu widerrufen. Andernfalls können sie diese Email ignorieren.</strong></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
				}
			},
			"confirmEmailChangeTemplate": {
				"body": "<p>Hallo,</p>\n<p>klicken sie bitte auf den folgenden Knopf, um ihre neue Email Addresse zu bestätigen.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Email Bestätigen</a>\n</p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			},
			"otp": {
				"emailTemplate": {
					"body": "<p>Hallo,</p>\n<p>bitte klicken auf den folgenden Knopf, um sich anzumelden:</p>\n  <a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende Einmalpasswort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn sie nicht versuchen, sich anzumelden können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
				}
			},
			"resetPasswordTemplate": {
				"body": "<p>Hallo,</p>\n<p>zum zurücksetzen ihres Passwort klicken sie btte auf den folgenden Knopf.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Passwort zurücksetzen</a>\n</p>\n<p><i>Wenn sie nicht versucht haben, ihr Passwort zurückzusetzen, können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			},
			"verificationTemplate": {
				"body": "<p>Hallo,</p>\n<p>bitte klicken sie auf den folgenden Knopf, um ihre Mail Addresse zu bestätigen.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Bestätigen</a>\n</p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>"
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
