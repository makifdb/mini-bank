package email

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/wneessen/go-mail"
)

// GoMailService is an implementation of EmailService using go-mail
type MailService struct {
	smtpServer string
	username   string
	password   string
	port       int
}

// NewMailService creates a new MailService
func NewMailService(smtpServer, username, password string, port int) *MailService {
	return &MailService{
		smtpServer: smtpServer,
		username:   username,
		password:   password,
		port:       port,
	}
}

// SendEmail sends an email using go-mail
func (s *MailService) SendEmail(to, subject, body string) error {
	m := mail.NewMsg()
	if err := m.From("info@mini-bank.com"); err != nil {
		return fmt.Errorf("failed to set from address: %w", err)
	}
	if err := m.To(to); err != nil {
		return fmt.Errorf("failed to set to address: %w", err)
	}
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, body)
	c, err := mail.NewClient(s.smtpServer, mail.WithPort(s.port), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.username), mail.WithPassword(s.password))
	if err != nil {
		return fmt.Errorf("failed to create email client: %w", err)
	}
	if err := c.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// SendVerificationEmail sends a verification email to the user
func (s *MailService) SendVerificationEmail(to, code string) error {
	return s.SendEmail(to, "Verification Code", "Your verification code is: "+code)
}

func (s *MailService) SendLoginEmail(to, code string) error {
	return s.SendEmail(to, "Login Code", "Your login code is: "+code)
}

// SendTransferNotification sends a transfer notification to both sender and receiver
func (s *MailService) SendTransferNotification(fromEmail, toEmail string, amount decimal.Decimal) error {
	fromBody := fmt.Sprintf("You have successfully transferred %s to %s", amount.String(), toEmail)
	toBody := fmt.Sprintf("You have received %s from %s", amount.String(), fromEmail)

	if err := s.SendEmail(fromEmail, "Transfer Notification", fromBody); err != nil {
		return fmt.Errorf("failed to send email to sender: %w", err)
	}

	if err := s.SendEmail(toEmail, "Transfer Notification", toBody); err != nil {
		return fmt.Errorf("failed to send email to receiver: %w", err)
	}

	return nil
}

// SendWithdrawalNotification sends a withdrawal notification to the user
func (s *MailService) SendWithdrawalNotification(toEmail string, amount decimal.Decimal) error {
	body := fmt.Sprintf("You have successfully withdrawn %s", amount.String())
	return s.SendEmail(toEmail, "Withdrawal Notification", body)
}

// SendDepositNotification sends a deposit notification to the user
func (s *MailService) SendDepositNotification(toEmail string, amount decimal.Decimal) error {
	body := fmt.Sprintf("You have successfully deposited %s", amount.String())
	return s.SendEmail(toEmail, "Deposit Notification", body)
}
