package dto

type SignUpRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerificationRequest struct {
	Email            string `json:"email" binding:"required,email"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
