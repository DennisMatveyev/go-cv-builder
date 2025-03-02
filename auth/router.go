package auth

import (
	"cv-builder/common"
	"cv-builder/models"
	"cv-builder/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(r fiber.Router, db *gorm.DB, jwtSecret string) {
	userRepo := user.NewUserRepository(db)
	authService := NewAuthService(userRepo, jwtSecret)

	r.Post("/register", func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := common.ValidateRequest(c, user); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		if err := authService.Register(user); err != nil {
			switch err {
			case ErrUserExists:
				return common.ErrorResponse(c, fiber.StatusConflict, err)
			case common.ErrDatabase, ErrSaveUser, ErrHashPassword:
				return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
			default:
				return common.ErrorResponse(c, fiber.StatusBadRequest, err)
			}
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User signed up successfully"})
	})

	r.Post("/login", func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, common.ErrParseJSON)
		}
		token, err := authService.Login(user.Email, user.Password)
		if err != nil {
			switch err {
			case ErrGenerateToken:
				return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
			default:
				return common.ErrorResponse(c, fiber.StatusBadRequest, err)
			}
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
	})
}
