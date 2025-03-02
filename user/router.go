package user

import (
	"cv-builder/common"
	"cv-builder/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(r fiber.Router, db *gorm.DB) {
	userRepo := NewUserRepository(db)

	r.Get("/profile", func(c *fiber.Ctx) error {
		profile, err := userRepo.GetProfile(c.Locals("userID").(int))
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(profile)
	})

	r.Get("/jobs", func(c *fiber.Ctx) error {
		jobs, err := userRepo.GetJobs(c.Locals("userID").(int))
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(jobs)
	})

	r.Post("/profile", func(c *fiber.Ctx) error {
		profile := new(models.Profile)
		if err := common.ValidateRequest(c, profile); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		profileID, err := userRepo.CreateProfile(c.Locals("userID").(int), profile)
		if err != nil {
			switch err {
			case ErrProfileExists:
				return common.ErrorResponse(c, fiber.StatusBadRequest, err)
			default:
				return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
			}
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":   "Profile created",
			"profileID": profileID,
		})
	})

	r.Post("/jobs", func(c *fiber.Ctx) error {
		job := new(models.Job)
		if err := common.ValidateRequest(c, job); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		jobID, err := userRepo.CreateJob(c.Locals("userID").(int), job)
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Jobs created",
			"jobID":   jobID,
		})
	})

	r.Put("/profile", func(c *fiber.Ctx) error {
		profile := new(models.Profile)
		if err := common.ValidateRequest(c, profile); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		err := userRepo.UpdateProfile(profile)
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profile updated"})
	})

	r.Put("/jobs", func(c *fiber.Ctx) error {
		job := new(models.Job)
		if err := common.ValidateRequest(c, job); err != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, err)
		}
		err := userRepo.UpdateJob(job)
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Job updated"})
	})

	r.Delete("/jobs/:id", func(c *fiber.Ctx) error {
		jobID, convErr := c.ParamsInt("id")
		if convErr != nil {
			return common.ErrorResponse(c, fiber.StatusBadRequest, ErrInvalidParamFormat)
		}
		err := userRepo.DeleteJob(jobID)
		if err != nil {
			return common.ErrorResponse(c, fiber.StatusInternalServerError, err)
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Job deleted"})
	})
}
