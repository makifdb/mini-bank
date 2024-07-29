package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/makifdb/mini-bank/corporate/internal/application/service"
	"github.com/makifdb/mini-bank/corporate/internal/interface/dto"
	"github.com/makifdb/mini-bank/corporate/pkg/error"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService *service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterRoutes(app *gin.RouterGroup) {
	app.POST("/users", h.createUser())
	app.GET("/users", h.getUsers())
	app.GET("/users/:id", h.getUser())
	app.PATCH("/users/:id", h.updateUser())
	app.DELETE("/users/:id", h.deleteUser())
}

// createUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  dto.CreateUserRequest  true  "User info"
// @Success 201 {object} user.User
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /users [post]
func (h *UserHandler) createUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
			return
		}

		user, err := h.userService.CreateUser(c, req.FirstName, req.LastName, req.Email)
		if err != nil {
			h.logger.Error("Failed to create user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to create user")))
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

// getUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept  json
// @Produce  json
// @Param   limit  query  int  false  "Limit"
// @Param   offset  query  int  false  "Offset"
// @Success 200 {object} []user.User
// @Failure 500 {object} error.APIError
// @Router /users [get]
func (h *UserHandler) getUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset := 10, 0
		if l, ok := c.GetQuery("limit"); ok {
			if lInt, err := strconv.Atoi(l); err == nil {
				limit = lInt
			}
		}
		if o, ok := c.GetQuery("offset"); ok {
			if oInt, err := strconv.Atoi(o); err == nil {
				offset = oInt
			}
		}

		users, err := h.userService.GetUsers(c, limit, offset)
		if err != nil {
			h.logger.Error("Failed to find users", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to find users")))
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// getUser godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "User ID"
// @Success 200 {object} user.User
// @Failure 400 {object} error.APIError
// @Failure 404 {object} error.APIError
// @Router /users/{id} [get]
func (h *UserHandler) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.Error("Failed to parse ID", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
			return
		}

		user, err := h.userService.GetUser(c, id)
		if err != nil {
			h.logger.Error("Failed to find user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusNotFound, error.NewErrorResponse(*error.NotFound("User not found")))
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// updateUser godoc
// @Summary Update user
// @Description Update user
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "User ID"
// @Param   user  body  dto.UpdateUserRequest  true  "User info"
// @Success 200 {object} user.User
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /users/{id} [patch]
func (h *UserHandler) updateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.logger.Error("Failed to bind JSON", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
			return
		}

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.Error("Failed to parse ID", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
			return
		}

		user, err := h.userService.UpdateUser(c, id, req.FirstName, req.LastName, req.Email)
		if err != nil {
			h.logger.Error("Failed to update user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to update user")))
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// deleteUser godoc
// @Summary Delete user
// @Description Delete user
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "User ID"
// @Success 204
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /users/{id} [delete]
func (h *UserHandler) deleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.Error("Failed to parse ID", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest("Invalid ID")))
			return
		}

		if err := h.userService.DeleteUser(c, id); err != nil {
			h.logger.Error("Failed to delete user", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError("Failed to delete user")))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}
