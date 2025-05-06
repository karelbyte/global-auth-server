package controllers

import (
	"net/http"
	"time"

	"global-auth-server/libs"
	"global-auth-server/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents the request body for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents the response body for user information.
type UserResponse struct {
	ID              string        `json:"id"`
	Username        string        `json:"username"`
	Code            *string       `json:"code,omitempty"`
	Names           string        `json:"names"`
	Email           string        `json:"email"`
	RolID           *string       `json:"rol_id,omitempty"`
	IsStaff         bool          `json:"is_staff"`
	IsActive        bool          `json:"is_active"`
	BossID          *string       `json:"boss_id,omitempty"`
	Logins          *int          `json:"logins,omitempty"`
	CanDownloadXlsx *bool         `json:"can_download_xlsx,omitempty"`
	BankID          *string       `json:"bank_id,omitempty"`
	FilialID        *string       `json:"filial_id,omitempty"`
	Roles           []services.Role `json:"roles"`
}

// LoginResponse represents the response body for a successful login.
type LoginResponse struct {
	Message   string       `json:"message"`
	User      UserResponse `json:"user"`
	Token     string       `json:"token"`
	ExpiredAt int64        `json:"expired_at"`
}

// ErrorResponse represents a generic error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Login godoc
// @Summary Authenticate user and return JWT token
// @Description Authenticate user with email and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "User credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	loggingService := services.NewLoggingService()
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	user, err := services.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user or password"})
		return
	}

	// Validate password using bcrypt
	if user.Password == nil || bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user or password"})
		return
	}

	// Get user roles
	roles, err := services.GetRolesByUserID(user.ID)
	if err != nil {
		roles = []services.Role{}
	}

	userResponse := UserResponse{
		ID:              user.ID,
		Username:        user.Username,
		Code:            user.Code,
		Names:           user.Names,
		Email:           user.Email,
		RolID:           user.RolID,
		IsStaff:         user.IsStaff,
		IsActive:        user.IsActive,
		BossID:          user.BossID,
		Logins:          user.Logins,
		CanDownloadXlsx: user.CanDownloadXlsx,
		BankID:          user.BankID,
		FilialID:        user.FilialID,
		Roles:           roles,
	}

	// JWT payload (puedes ajustar los campos que quieras incluir)
	payload := map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   roles,
	}

	token, expiredAt, err := libs.GenerateJWT(payload, 30*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	loginResponse := LoginResponse{
		Message:   "Login successful",
		User:      userResponse,
		Token:     token,
		ExpiredAt: expiredAt,
	}

	loggingService.Log(user.ID, c.Request.URL.Path, nil, gin.H{"message": "Login successful", "user": userResponse, "token": token, "expired_at": expiredAt}, "LOGIN_SUCCESS")

	c.JSON(http.StatusOK, loginResponse)
}

// Home godoc
// @Summary Show the home page
// @Description Get the home page message
// @Tags home
// @Produce html
// @Success 200 {string} string "OK"
// @Router / [get]
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"message": "Todas las peticiones se atienden con el prefijo /api",
	})
}