package admin

import (
	"fmt"
	"regexp"

	"github.com/makifdb/mini-bank/corporate/internal/domain/base"
)

// Admin represents an admin user.
type Admin struct {
	base.Base        // Embedded Base struct to include common fields
	Email     string `json:"email" gorm:"unique"` // Email address of the admin
}

// NewUser creates a new User instance with validation.
func NewAdmin(email string) (*Admin, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email format")
	}
	return &Admin{
		Base:  base.NewBase(),
		Email: email,
	}, nil
}

// isValidEmail checks if the provided email address has a valid format.
func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
