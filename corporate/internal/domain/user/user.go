package user

import (
	"fmt"
	"regexp"

	"github.com/makifdb/mini-bank/corporate/internal/domain/base"
)

// User represents a bank account holder.
type User struct {
	base.Base        // Embedded Base struct to include common fields
	FirstName string `json:"first_name" gorm:"index"`  // First name of the user
	LastName  string `json:"last_name" gorm:"index"`   // Last name of the user
	Email     string `json:"email" gorm:"uniqueIndex"` // Email address of the user
}

// NewUser creates a new User instance with validation.
func NewUser(firstName, lastName, email string) (*User, error) {

	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email address")
	}

	return &User{
		Base:      base.NewBase(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}

// FullName returns the full name of the user.
func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// isValidEmail checks if the provided email address has a valid format.
func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
