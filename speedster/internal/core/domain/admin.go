package domain

import "fmt"

type Admin struct {
	Base
	Email string `json:"email" pg:",unique,notnull"`
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
