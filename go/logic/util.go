package logic

import (
	"fmt"
	"os"
	"strconv"
	"time"

	mail "github.com/go-mail/mail/v2"
)

type BrevoConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
	FromName string
}

func LoadBrevoConfig() (BrevoConfig, error) {
	port, err := strconv.Atoi(os.Getenv("BREVO_PORT"))
	if err != nil {
		return BrevoConfig{}, fmt.Errorf("invalid port: %v", err)
	}

	config := BrevoConfig{
		Host:     os.Getenv("BREVO_HOST"),
		Port:     port,
		User:     os.Getenv("BREVO_LOGIN"),
		Pass:     os.Getenv("BREVO_PASS"),
		From:     os.Getenv("BREVO_FROM"),
		FromName: os.Getenv("BREVO_FROM_NAME"),
	}

	return config, nil
}

func SendEmailBrevo(config *BrevoConfig, to string, subject string, textBody string, htmlBody string) error {
	message := mail.NewMessage()

	if config.FromName != "" {
		message.SetAddressHeader("From", config.From, config.FromName)
	} else {
		message.SetHeader("From", config.From)
	}
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	if textBody == "" {
		textBody = "Please view this email in an HTML-compatible email viewer."
	}

	message.SetBody("text/plain", textBody)
	if htmlBody != "" {
		message.AddAlternative("text/html", htmlBody)
	}

	dialer := mail.NewDialer(config.Host, config.Port, config.User, config.Pass)
	dialer.StartTLSPolicy = mail.MandatoryStartTLS
	dialer.Timeout = 10 * time.Second

	return dialer.DialAndSend(message)
}
