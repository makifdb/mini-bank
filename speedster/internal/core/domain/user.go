package domain

import (
	"fmt"
	"regexp"
)

// User represents a bank account holder.
type User struct {
	Base
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Accounts  []Account `json:"accounts" pg:"rel:has-many"`
}

func NewUser(firstName, lastName, email string) (*User, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email address")
	}
	return &User{
		Base:      NewBase(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
