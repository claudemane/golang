package v1

import (
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, t usecase.UserInterface, l logger.Interface) {
	// Global rate limiter: 60 requests per minute per identity.
	// Authenticated users are keyed by their JWT UserID.
	// Anonymous requests are keyed by their client IP.
	rateLimiter := utils.NewRateLimiter(5, time.Minute)
	handler.Use(rateLimiter.Middleware())

	v1 := handler.Group("/v1")
	{
		newUserRoutes(v1, t, l)
	}
}
