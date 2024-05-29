package MailService

import (
	"awesomeProject/model"
)
import "gopkg.in/gomail.v2"

func SendMail(account model.EMailAccountConfig, tarEmail string, message string, subject string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "CyberAster"+"<"+account.Account+">")
	m.SetHeader("To", tarEmail)
	m.SetHeader("subject", "Hello!")
	m.SetBody("text/html", message)
	dialer := gomail.NewDialer(
		model.EMail.Server,
		model.EMail.Portal,
		account.Account,
		account.Password)
	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
