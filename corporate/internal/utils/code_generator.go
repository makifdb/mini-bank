package utils

import (
	"crypto/rand"
)

const charset = "0123456789"
const length = 6

func GenerateVerificationCode() (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b), nil
}
