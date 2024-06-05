package implementation

import (
	"context"
	"fmt"
	"net/smtp"
	"sso-service/config"
	"sso-service/internal/domain/service"
)

type MailService struct {
	config config.SmtpConfig
	logger service.Logger
}

func NewMailService(config config.SmtpConfig, logger service.Logger) service.MailService {
	return &MailService{
		config: config,
		logger: logger,
	}
}

func (ms *MailService) SendMail(ctx context.Context, address, header, message string) error {
	auth := smtp.PlainAuth("", ms.config.SmtpFrom, ms.config.SmtpPassword, ms.config.SmtpHost)
	finalMessage := fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		address,
		header,
		message)

	err := smtp.SendMail(ms.config.SmtpHost+":"+ms.config.SmtpPort,
		auth,
		ms.config.SmtpFrom,
		[]string{
			address,
		},
		[]byte(finalMessage))

	if err != nil {
		ms.logger.Error(fmt.Sprintf("[MailService] ошибка отправки сообщения на %s : %s", address, err.Error()))
		return service.ErrorMailSending
	}

	return nil
}
