package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/makifdb/mini-bank/corporate/internal/application/service"
	"github.com/makifdb/mini-bank/corporate/internal/interface/dto"
	"github.com/makifdb/mini-bank/corporate/pkg/error"
	"go.uber.org/zap"
)

type AdminHandler struct {
	adminService *service.AdminService
	logger       *zap.Logger
}

func NewAdminHandler(adminService *service.AdminService, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		logger:       logger,
	}
}

func (h *AdminHandler) RegisterRoutes(app *gin.RouterGroup) {
	app.POST("/signup", h.signUp)
	app.POST("/login", h.login)
	app.POST("/verify", h.verifyAdmin)
	app.POST("/refresh", h.refreshToken)
	app.POST("/logout", h.logout)
}

// signUp godoc
// @Summary Sign up as an admin
// @Description Sign up as an admin
// @Tags admins
// @Accept  json
// @Produce  json
// @Param   admin  body  dto.SignUpRequest  true  "Admin info"
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /admins/signup [post]
func (h *AdminHandler) signUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
		return
	}

	admin, err := h.adminService.SignUp(c, req.Email)
	if err != nil {
		h.logger.Error("Failed to create admin", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError(err.Error())))
		return
	}

	h.logger.Info("Admin created successfully", zap.String("email", admin.Email))
	c.JSON(http.StatusCreated, gin.H{"message": "Admin created successfully, check your email for verification code"})
}

// login godoc
// @Summary Login as an admin
// @Description Login as an admin
// @Tags admins
// @Accept  json
// @Produce  json
// @Param   admin  body  dto.LoginRequest  true  "Admin info"
// @Failure 400 {object} error.APIError
// @Failure 404 {object} error.APIError
// @Router /admins/login [post]
func (h *AdminHandler) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
		return
	}

	err := h.adminService.Login(c, req.Email)
	if err != nil {
		h.logger.Error("Admin not found", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusNotFound, error.NewErrorResponse(*error.NotFound("Admin not found")))
		return
	}

	h.logger.Info("Login code sent successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, gin.H{"message": "Login code sent successfully"})
}

// verifyAdmin godoc
// @Summary Verify an admin
// @Description Verify an admin
// @Tags admins
// @Accept  json
// @Produce  json
// @Param   verification  body  dto.VerificationRequest  true  "Verification info"
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /admins/verify [post]
func (h *AdminHandler) verifyAdmin(c *gin.Context) {
	var req dto.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
		return
	}

	tokenString, err := h.adminService.VerifyAdmin(c, req.Email, req.VerificationCode)
	if err != nil {
		h.logger.Error("Verification failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError(err.Error())))
		return
	}

	h.logger.Info("Admin verified and token generated", zap.String("email", req.Email))
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// refreshToken godoc
// @Summary Refresh an admin token
// @Description Refresh an admin token
// @Tags admins
// @Accept  json
// @Produce  json
// @Param   token  body  dto.RefreshTokenRequest  true  "Token info"
// @Failure 400 {object} error.APIError
// @Failure 500 {object} error.APIError
// @Router /admins/refresh [post]
func (h *AdminHandler) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, error.NewErrorResponse(*error.BadRequest(err.Error())))
		return
	}

	newTokenString, err := h.adminService.RefreshToken(c, req.Token)
	if err != nil {
		h.logger.Error("Failed to refresh token", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, error.NewErrorResponse(*error.InternalServerError(err.Error())))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}

// logout godoc
// @Summary Logout an admin
// @Description Logout an admin
// @Tags admins
// @Router /admins/logout [post]
func (h *AdminHandler) logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
