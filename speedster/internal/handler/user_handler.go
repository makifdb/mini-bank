package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/makifdb/mini-bank/speedster/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/users", h.CreateUser)
	api.Get("/users/:id", h.GetUserByID)
	api.Put("/users/:id", h.UpdateUser)
	api.Delete("/users/:id", h.DeleteUser)
	api.Get("/users", h.ListUsers)
	api.Get("/users/:id/accounts", h.GetUserByIDWithAccounts)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "Create user"
// @Success 200 {object} models.User
// @Failure 400 {object} fiber.Map
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	type request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.userService.CreateUser(c.Context(), req.FirstName, req.LastName, req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(user)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} fiber.Map
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(http.StatusOK).JSON(user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's details
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param user body models.User true "Update user"
// @Success 200 {object} models.User
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	type request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.userService.UpdateUser(c.Context(), id, req.FirstName, req.LastName, req.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Param id path string true "User ID"
// @Success 204
// @Failure 404 {object} fiber.Map
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.userService.DeleteUser(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(http.StatusNoContent).Send(nil)
}

// ListUsers godoc
// @Summary List all users
// @Description List all users with pagination
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.User
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	users, err := h.userService.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(users)
}

// GetUserByIDWithAccounts godoc
// @Summary Get a user by ID with accounts
// @Description Get a user by ID with associated accounts
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} fiber.Map
// @Router /users/{id}/accounts [get]
func (h *UserHandler) GetUserByIDWithAccounts(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUserByIDWithAccounts(c.Context(), id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(http.StatusOK).JSON(user)
}
