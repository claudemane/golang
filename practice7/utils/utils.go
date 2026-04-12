package utils

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ---------------------------------------------------------------------------
// Password helpers
// ---------------------------------------------------------------------------

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// ---------------------------------------------------------------------------
// JWT helpers
// ---------------------------------------------------------------------------

func GenerateJWT(userID uuid.UUID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ---------------------------------------------------------------------------
// Middleware: JWTAuthMiddleware
// Validates the Bearer token, then sets "userID" and "role" in the context.
// ---------------------------------------------------------------------------

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		c.Set("userID", claims["user_id"].(string))
		c.Set("role", claims["role"].(string))
		c.Next()
	}
}

// ---------------------------------------------------------------------------
// Middleware: RoleMiddleware
// Must be used AFTER JWTAuthMiddleware (needs "role" already in context).
// Returns 403 Forbidden if the caller's role doesn't match requiredRole.
// ---------------------------------------------------------------------------

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			return
		}
		if role.(string) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied: requires role '" + requiredRole + "'",
			})
			return
		}
		c.Next()
	}
}

// ---------------------------------------------------------------------------
// Middleware: RateLimiterMiddleware
//
// Identification strategy:
//   - Authenticated request  → key = "user:<userID>"   (from JWT claim)
//   - Anonymous request      → key = "ip:<ClientIP>"
//
// Each key gets a fixed window of `maxRequests` per `window` duration.
// Exceeding the limit returns HTTP 429.
// sync.Mutex prevents race conditions on the shared counters map.
// ---------------------------------------------------------------------------

type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	counters    map[string]*rateLimitEntry
	maxRequests int
	window      time.Duration
}

func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		counters:    make(map[string]*rateLimitEntry),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Determine the key: prefer JWT userID, fall back to client IP.
		key := "ip:" + c.ClientIP()
		if userID, exists := c.Get("userID"); exists && userID != "" {
			key = "user:" + userID.(string)
		} else {
			// Try to extract user_id from the Authorization header even when
			// JWTAuthMiddleware is not in the chain (rate-limit on public routes).
			tokenStr := c.GetHeader("Authorization")
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			if tokenStr != "" {
				if token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				}); err == nil && token.Valid {
					if claims, ok := token.Claims.(jwt.MapClaims); ok {
						if uid, ok := claims["user_id"].(string); ok && uid != "" {
							key = "user:" + uid
						}
					}
				}
			}
		}

		rl.mu.Lock()
		entry, exists := rl.counters[key]
		now := time.Now()

		if !exists || now.After(entry.windowEnd) {
			// Start a fresh window.
			rl.counters[key] = &rateLimitEntry{count: 1, windowEnd: now.Add(rl.window)}
			rl.mu.Unlock()
			c.Next()
			return
		}

		entry.count++
		if entry.count > rl.maxRequests {
			rl.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}
		rl.mu.Unlock()
		c.Next()
	}
}
