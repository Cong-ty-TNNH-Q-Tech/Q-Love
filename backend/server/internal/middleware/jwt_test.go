package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware(t *testing.T) {
	app := fiber.New()
	
	// Create route to test middleware
	app.Get("/test", JWTMiddleware("test-secret"), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		return c.SendString(userID.String())
	})
	
	// Also test default secret
	app.Get("/default-secret", JWTMiddleware(""), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(uuid.UUID)
		return c.SendString(userID.String())
	})

	validUUID := uuid.New()

	generateToken := func(secret string, claims jwt.MapClaims, alg jwt.SigningMethod) string {
		token := jwt.NewWithClaims(alg, claims)
		tokenString, _ := token.SignedString([]byte(secret))
		return tokenString
	}

	tests := []struct {
		name           string
		route          string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Missing Authorization header",
			route:          "/test",
			authHeader:     "",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization format (no Bearer)",
			route:          "/test",
			authHeader:     "Basic token123",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Invalid Authorization format (only Bearer)",
			route:          "/test",
			authHeader:     "Bearer",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Invalid token signature",
			route:          "/test",
			authHeader:     "Bearer " + generateToken("wrong-secret", jwt.MapClaims{"sub": validUUID.String()}, jwt.SigningMethodHS256),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Invalid algorithm",
			route:          "/test",
			// We simulate none alg vulnerability
			authHeader:     "Bearer " + generateToken("test-secret", jwt.MapClaims{"sub": validUUID.String()}, jwt.SigningMethodNone),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Expired token",
			route:          "/test",
			authHeader:     "Bearer " + generateToken("test-secret", jwt.MapClaims{"sub": validUUID.String(), "exp": time.Now().Add(-1 * time.Hour).Unix()}, jwt.SigningMethodHS256),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Invalid UUID subject",
			route:          "/test",
			authHeader:     "Bearer " + generateToken("test-secret", jwt.MapClaims{"sub": "invalid-uuid"}, jwt.SigningMethodHS256),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Missing sub claim",
			route:          "/test",
			authHeader:     "Bearer " + generateToken("test-secret", jwt.MapClaims{"role": "user"}, jwt.SigningMethodHS256),
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "Valid token",
			route:          "/test",
			authHeader:     "Bearer " + generateToken("test-secret", jwt.MapClaims{"sub": validUUID.String()}, jwt.SigningMethodHS256),
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Valid token with default secret",
			route:          "/default-secret",
			authHeader:     "Bearer " + generateToken(DefaultSecret, jwt.MapClaims{"sub": validUUID.String()}, jwt.SigningMethodHS256),
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.route, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			resp, _ := app.Test(req)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
