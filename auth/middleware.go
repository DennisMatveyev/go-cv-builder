package auth

import (
	"cv-builder/common"
	"cv-builder/models"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func Middleware(db *gorm.DB, jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return common.ErrorResponse(c, fiber.StatusUnauthorized, ErrMissingAuthHeader)
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return common.ErrorResponse(c, fiber.StatusUnauthorized, err)
		}

		userID := token.Claims.(jwt.MapClaims)["user_id"]
		var user models.User
		err = db.First(&user, userID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrorResponse(c, fiber.StatusUnauthorized, ErrInvalidToken)
		}
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, common.ErrDatabase)
		}
		c.Locals("userID", int(userID.(float64)))

		return c.Next()
	}
}
