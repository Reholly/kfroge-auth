package service

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrorMailSending = errors.New("ошибка отправки сообщения на почту")
)

type MailService interface {
	SendMail(ctx context.Context, address, header, message string) error
}
