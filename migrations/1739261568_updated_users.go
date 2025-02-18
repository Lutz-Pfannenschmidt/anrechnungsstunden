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
					"body": "<p>Hallo,</p>\n<p>Wir haben festgestellt, dass Sie sich von einem neuen Standort aus angemeldet haben</p>\n<p><strong>Wenn Sie das nicht waren, sollten Sie sofort das Passwort für Ihr Konto ändern, um den Zugriff von allen anderen Standorten zu widerrufen. Andernfalls können sie diese Email ignorieren.</strong></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>",
					"subject": "{APP_NAME} anmeldung"
				}
			},
			"confirmEmailChangeTemplate": {
				"body": "<p>Hallo,</p>\n<p>klicken sie bitte auf den folgenden Knopf, um ihre neue Email Addresse zu bestätigen.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Email Bestätigen</a>\n</p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>",
				"subject": "Neue Email Addresse für {APP_NAME} bestätigen"
			},
			"fileToken": {
				"duration": 1800
			},
			"otp": {
				"emailTemplate": {
					"body": "<p>Hallo,</p>\n<p>bitte klicken auf den folgenden Knopf, um sich anzumelden:</p>\n  <a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende Einmalpasswort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn sie nicht versuchen, sich anzumelden können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>",
					"subject": "Einmalpasswort für {APP_NAME}"
				}
			},
			"resetPasswordTemplate": {
				"body": "<p>Hallo,</p>\n<p>zum zurücksetzen ihres Passwort klicken sie btte auf den folgenden Knopf.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>Wenn sie nicht versucht haben, ihr Passwort zurückzusetzen, können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>",
				"subject": "Passwort für {APP_NAME} zurücksetzen"
			},
			"verificationTemplate": {
				"body": "<p>Hallo,</p>\n<p>bitte klicken sie auf den folgenden Knopf, um ihre Mail Addresse zu bestätigen.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Bestätigen</a>\n</p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>",
				"subject": "Email Verifizierung für {APP_NAME}"
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
					"body": "<p>Hello,</p>\n<p>We noticed a login to your {APP_NAME} account from a new location.</p>\n<p>If this was you, you may disregard this email.</p>\n<p><strong>If this wasn't you, you should immediately change your {APP_NAME} account password to revoke access from all other locations.</strong></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
					"subject": "Login from a new location"
				}
			},
			"confirmEmailChangeTemplate": {
				"body": "<p>Hello,</p>\n<p>Click on the button below to confirm your new email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-email-change/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Confirm new email</a>\n</p>\n<p><i>If you didn't ask to change your email address, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Confirm your {APP_NAME} new email address"
			},
			"fileToken": {
				"duration": 180
			},
			"otp": {
				"emailTemplate": {
					"body": "<p>Hallo,</p>\n<p>bitte klicken auf den folgenden Knopf, um sich anzumelden:</p>\n  <a class=\"btn\" href=\"{APP_URL}/otp/?id={OTP_ID}&otp={OTP}\" target=\"_blank\" rel=\"noopener\">Anmelden</a>\n<p>Alternativ können Sie auch das folgende einmal Passwort verwenden: <strong>{OTP}</strong></p>\n<p>Bitte kopieren Sie diesen Code und geben Sie ihn auf der Anmeldeseite ein, um sich anzumelden.</p>\n<p><i>Wenn sie nicht versuchen, sich anzumelden können sie diese Email ignorieren.</i></p>\n<p>\n  Diese E-Mail wurde automatisch generiert und Antworten auf diese werden nicht gelesen. Bei Fragen wenden sie sich bitte an die für dieses Programm zuständige Person.\n</p>",
					"subject": "OTP for {APP_NAME}"
				}
			},
			"resetPasswordTemplate": {
				"body": "<p>Hello,</p>\n<p>Click on the button below to reset your password.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-password-reset/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Reset password</a>\n</p>\n<p><i>If you didn't ask to reset your password, you can ignore this email.</i></p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Reset your {APP_NAME} password"
			},
			"verificationTemplate": {
				"body": "<p>Hello,</p>\n<p>Thank you for joining us at {APP_NAME}.</p>\n<p>Click on the button below to verify your email address.</p>\n<p>\n  <a class=\"btn\" href=\"{APP_URL}/_/#/auth/confirm-verification/{TOKEN}\" target=\"_blank\" rel=\"noopener\">Verify</a>\n</p>\n<p>\n  Thanks,<br/>\n  {APP_NAME} team\n</p>",
				"subject": "Verify your {APP_NAME} email"
			}
		}`), &collection); err != nil {
			return err
		}

		return app.Save(collection)
	})
}
