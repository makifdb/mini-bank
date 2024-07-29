package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/makifdb/mini-bank/speedster/internal/service"
)

type AuthHandler struct {
	adminService *service.AdminService
}

func NewAuthHandler(adminService *service.AdminService) *AuthHandler {
	return &AuthHandler{
		adminService: adminService,
	}
}

func (h *AuthHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/signup", h.signUp)
	app.Post("/login", h.login)
	app.Post("/verify", h.verifyAdmin)
	app.Post("/refresh", h.refreshToken)
	app.Post("/logout", h.logout)
}

// signUp godoc
// @Summary Create a new admin account
// @Description Create a new admin account
// @Tags auth
// @Accept json
// @Produce json
// @Param admin body SignUpRequest true "Admin request"
// @Success 200 {object} SignUpResponse
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /signup [post]
func (h *AuthHandler) signUp(c *fiber.Ctx) error {

	type SignUpRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	admin, err := h.adminService.SignUp(c.Context(), req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"email": admin.Email})
}

func (h *AuthHandler) login(c *fiber.Ctx) error {

	type LoginRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := h.adminService.Login(c.Context(), req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(http.StatusOK)
}

func (h *AuthHandler) verifyAdmin(c *fiber.Ctx) error {
	type VerifyRequest struct {
		Email string `json:"email" validate:"required,email"`
		Code  string `json:"code" validate:"required"`
	}

	var req VerifyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := h.adminService.VerifyAdmin(c.Context(), req.Email, req.Code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) refreshToken(c *fiber.Ctx) error {
	var RefreshRequest struct {
		Token string `json:"token" validate:"required"`
	}

	if err := c.BodyParser(&RefreshRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	token, err := h.adminService.RefreshToken(c.Context(), RefreshRequest.Token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) logout(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
