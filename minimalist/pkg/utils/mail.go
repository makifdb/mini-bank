package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type MailService struct {
	smtpServer string
	username   string
	password   string
	port       int
	auth       smtp.Auth
}

func NewMailService() *MailService {
	p, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_SERVER"))

	return &MailService{
		smtpServer: os.Getenv("SMTP_SERVER"),
		username:   os.Getenv("SMTP_USER"),
		password:   os.Getenv("SMTP_PASSWORD"),
		port:       p,
		auth:       auth,
	}
}

func (s *MailService) SendEmail(to, subject, body string) error {
	from := "info@mini-bank.com"
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := fmt.Sprintf("%s:%d", s.smtpServer, s.port)

	if err := smtp.SendMail(addr, s.auth, from, strings.Split(to, ","), msg); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *MailService) SendVerificationEmail(to, code string) error {
	return s.SendEmail(to, "Verification Code", "Your verification code is: "+code)
}

func (s *MailService) SendLoginEmail(to, code string) error {
	return s.SendEmail(to, "Login Code", "Your login code is: "+code)
}

func (s *MailService) SendTransferNotification(fromEmail, toEmail string, amount string) error {
	fromBody := fmt.Sprintf("You have successfully transferred %s to %s", amount, toEmail)
	toBody := fmt.Sprintf("You have received %s from %s", amount, fromEmail)

	if err := s.SendEmail(fromEmail, "Transfer Notification", fromBody); err != nil {
		return fmt.Errorf("failed to send email to sender: %w", err)
	}

	if err := s.SendEmail(toEmail, "Transfer Notification", toBody); err != nil {
		return fmt.Errorf("failed to send email to receiver: %w", err)
	}

	return nil
}

func (s *MailService) SendWithdrawalNotification(toEmail string, amount string) error {
	body := fmt.Sprintf("You have successfully withdrawn %s", amount)
	return s.SendEmail(toEmail, "Withdrawal Notification", body)
}

func (s *MailService) SendDepositNotification(toEmail string, amount string) error {
	body := fmt.Sprintf("You have successfully deposited %s", amount)
	return s.SendEmail(toEmail, "Deposit Notification", body)
}
