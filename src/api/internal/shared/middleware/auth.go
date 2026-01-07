package middleware

import (
	"strings"

	"github.com/algo-shield/algo-shield/src/api/internal/auth"
	apierrors "github.com/algo-shield/algo-shield/src/pkg/errors"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(authHandler *auth.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "Authorization header required"))
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "Invalid authorization header format"))
		}

		token := parts[1]
		user, err := authHandler.ValidateToken(token)
		if err != nil {
			// Check if it's an APIError with specific error code
			if apiErr, ok := err.(*apierrors.APIError); ok {
				return apierrors.SendError(c, apiErr)
			}
			// Fallback to generic unauthorized error
			return apierrors.SendError(c, apierrors.TokenInvalid())
		}

		// Store user in context
		c.Locals("user", user)
		c.Locals("user_id", user.ID)

		return c.Next()
	}
}

func RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok {
			return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "User not found in context"))
		}

		// Check if user has the required role
		hasRole := false
		for _, role := range user.Roles {
			if role.Name == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return apierrors.SendError(c, apierrors.InsufficientPermissions())
		}

		return c.Next()
	}
}

func RequireAnyRole(roleNames ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*models.User)
		if !ok {
			return apierrors.SendError(c, apierrors.NewAPIError(apierrors.ErrUnauthorized, "User not found in context"))
		}

		// Check if user has any of the required roles
		hasRole := false
		roleMap := make(map[string]bool)
		for _, roleName := range roleNames {
			roleMap[roleName] = true
		}

		for _, role := range user.Roles {
			if roleMap[role.Name] {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return apierrors.SendError(c, apierrors.InsufficientPermissions())
		}

		return c.Next()
	}
}
