package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	t usecase.UserInterface
	l logger.Interface
}

// NewUserRoutes wires all user-related routes.
//
// Public:
//   POST /v1/users/          → RegisterUser
//   POST /v1/users/login     → LoginUser
//
// Authenticated (JWT required):
//   GET  /v1/users/me                → GetMe
//   GET  /v1/users/protected/hello   → ProtectedFunc
//
// Admin only (JWT + role=admin required):
//   PATCH /v1/users/promote/:id      → PromoteUser
//
// Rate limiting is applied globally via the router group passed in.
func newUserRoutes(handler *gin.RouterGroup, t usecase.UserInterface, l logger.Interface) {
	r := &userRoutes{t, l}

	h := handler.Group("/users")
	{
		// Public routes
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		// Routes that require a valid JWT
		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/me", r.GetMe)
			protected.GET("/protected/hello", r.ProtectedFunc)

			// Admin-only sub-group
			adminOnly := protected.Group("/")
			adminOnly.Use(utils.RoleMiddleware("admin"))
			{
				adminOnly.PATCH("/promote/:id", r.PromoteUser)
			}
		}
	}
}

// ─── Handlers ───────────────────────────────────────────────────────────────

// RegisterUser creates a new user account with a hashed password.
func (r *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	createdUser, sessionID, err := r.t.RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully. Please check your email for verification code.",
		"session_id": sessionID,
		"user":       createdUser,
	})
}

// LoginUser validates credentials and returns a signed JWT.
func (r *userRoutes) LoginUser(c *gin.Context) {
	var input entity.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.t.LoginUser(&input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetMe returns the profile of the currently authenticated user.
//
// Requirements met:
//   - Works through JWTAuthMiddleware (userID stored in context).
//   - No body/query params needed — only the Authorization header.
//   - Calling with User A's token returns User A's data; User B's token → User B.
func (r *userRoutes) GetMe(c *gin.Context) {
	// JWTAuthMiddleware already validated the token and stored userID.
	rawID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userID, err := uuid.Parse(rawID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID in token"})
		return
	}

	user, err := r.t.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return only non-sensitive fields.
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"verified": user.Verified,
	})
}

// PromoteUser changes the target user's role to "admin".
//
// Requirements met:
//   - Protected by JWTAuthMiddleware + RoleMiddleware("admin").
//   - PATCH /v1/users/promote/:id
func (r *userRoutes) PromoteUser(c *gin.Context) {
	rawID := c.Param("id")
	targetID, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target user ID"})
		return
	}

	updatedUser, err := r.t.PromoteUser(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user": gin.H{
			"id":   updatedUser.ID,
			"username": updatedUser.Username,
			"role": updatedUser.Role,
		},
	})
}

// ProtectedFunc is a minimal demo of a JWT-protected endpoint.
func (r *userRoutes) ProtectedFunc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
