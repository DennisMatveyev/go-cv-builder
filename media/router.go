package media

import (
	"cv-builder/common"
	"cv-builder/config"
	"cv-builder/user"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(r fiber.Router, db *gorm.DB, cfg *config.Config) {
	userRepo := user.NewUserRepository(db)
	mediaService := NewMediaService(userRepo, cfg)

	r.Post("/upload_photo", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		if err := mediaService.UploadUserPhoto(c, file); err != nil {
			switch err {
			case ErrSaveFile, ErrRemoveOldFile:
				return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
			default:
				return common.ErrorResponse(c, fiber.StatusBadRequest, err)
			}
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "File uploaded successfully"})
	})

	r.Get("/download_cv", func(c *fiber.Ctx) error {
		filename, err := mediaService.DownloadCV(c)
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Download(filepath.Join(cfg.DownloadsDir, filename), "cv.pdf")
	})
}
