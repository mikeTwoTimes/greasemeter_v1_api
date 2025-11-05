package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer struct {
	client *sendgrid.Client
}

func NewMailer(client *sendgrid.Client) *Mailer {
	return &Mailer{client: client}
}

func (m *Mailer) SendPasswordReset(token string, email string) (int, gin.H) {
    from := mail.NewEmail("GreaseMeter", "no-reply@api.greasemeter.live")
    subject := "Reset Your Password"
    to := mail.NewEmail("", email)

    resetLink := fmt.Sprintf(
		"%s/reset-password/%s",
		"https://www.greasemeter.live/v1/users",
		token,
	)

    plainTextContent := fmt.Sprintf(
		"Click the following link to reset your password: %s",
		resetLink,
	)

	htmlContent := fmt.Sprintf(`
        <p>Click the following link to reset your password:</p>
        <p><a href="%s">%s</a></p>
        <p>If you didn’t request this, you can ignore this email.</p>
    `,
		resetLink,
		resetLink,
	)

    message := mail.NewSingleEmail(
		from,
		subject,
		to,
		plainTextContent,
		htmlContent,
	)

	response, err := m.client.Send(message)

	if err != nil {
        return http.StatusInternalServerError, gin.H{"error": err.Error()}
    } else if response.StatusCode >= 400 {
		return http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		}
	}

    return http.StatusNoContent, nil
}
