package models

import (
	"fmt"
)

type Admin struct {
	Base
	Email string `json:"email"` // Email address of the admin
}

func NewAdmin(email string) (*Admin, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email address")
	}

	return &Admin{
		Base:  NewBase(),
		Email: email,
	}, nil
}
