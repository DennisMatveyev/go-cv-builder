package media

import (
	"cv-builder/config"
	"cv-builder/models"
	"cv-builder/user"

	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

type MediaUploader interface {
	UploadUserPhoto(c *fiber.Ctx, file *multipart.FileHeader) error
}

type MediaDownloader interface {
	DownloadCV(c *fiber.Ctx) (string, error)
}

type MediaService interface {
	MediaUploader
	MediaDownloader
}

type mediaService struct {
	userRepo     user.UserRepository
	uploadsDir   string
	downloadsDir string
}

func NewMediaService(userRepo user.UserRepository, cfg *config.Config) MediaService {
	return &mediaService{
		userRepo:     userRepo,
		uploadsDir:   cfg.UploadsDir,
		downloadsDir: cfg.DownloadsDir,
	}
}

func (ms *mediaService) DownloadCV(c *fiber.Ctx) (string, error) {
	userID := c.Locals("userID").(int)
	_, err := os.Stat(ms.downloadsDir + fmt.Sprintf("%d.pdf", userID))
	if err == nil || !os.IsNotExist(err) {
		return fmt.Sprintf("%d.pdf", userID), nil
	}
	user, err := ms.userRepo.GetUserData(userID)
	if err != nil {
		return "", err
	}
	filename, err := ms.generateCV(user)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (ms *mediaService) generateCV(user *models.User) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	// Photo
	if user.Profile.PhotoPath != "" {
		_, err := os.Stat(user.Profile.PhotoPath)
		if err == nil {
			pdf.Image(user.Profile.PhotoPath, 10, 10, 40, 40, false, "", 0, "")
		}
	}
	pdf.SetY(55)
	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Curriculum Vitae")
	pdf.Ln(12)
	// Profile Section
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Name: %s %s", user.Profile.FirstName, user.Profile.LastName))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Email: %s", user.Email))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Phone: %s", user.Profile.PhoneNumber))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Header: %s", user.Profile.Header))
	pdf.Ln(15)
	// Jobs Section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Work Experience:")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)

	for _, job := range user.Jobs {
		pdf.Cell(40, 10, fmt.Sprintf("Company: %s", job.CompanyName))
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("Position: %s", job.Position))
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("Duration: %s - %s", job.StartDate, job.EndDate))
		pdf.Ln(8)
		pdf.MultiCell(0, 10, fmt.Sprintf("Description: %s", job.Description), "", "", false)
		pdf.Ln(10)
	}

	filename := fmt.Sprintf("%d.pdf", user.ID)
	err := pdf.OutputFileAndClose(ms.downloadsDir + filename)
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (ms *mediaService) UploadUserPhoto(c *fiber.Ctx, file *multipart.FileHeader) error {
	if err := ms.validateFileUpload(file.Size, file.Filename); err != nil {
		return err
	}
	if err := ms.removeOldUserPhoto(ms.uploadsDir, c.Locals("userID").(int)); err != nil {
		log.Println("Failed to remove old photo:", err)
		return ErrRemoveOldFile
	}

	newFileName := fmt.Sprintf("%d%s", c.Locals("userID").(int), filepath.Ext(file.Filename))
	savePath := filepath.Join(ms.uploadsDir, newFileName)

	if err := c.SaveFile(file, savePath); err != nil {
		log.Println("Failed to save photo:", err)
		return ErrSaveFile
	} else {
		userProfile, _ := ms.userRepo.GetProfile(c.Locals("userID").(int))
		userProfile.PhotoPath = savePath
		if err := ms.userRepo.UpdateProfile(userProfile); err != nil {
			log.Println("Failed to update profile when saving photo:", err)
			os.Remove(savePath)
			return ErrSaveFile
		}
	}
	return nil
}

func (ms *mediaService) validateFileUpload(fileSize int64, fileName string) error {
	if fileSize > 5*1024*1024 {
		return ErrFileTooLarge
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	if !slices.Contains([]string{".jpg", ".jpeg"}, ext) {
		return ErrInvalidFileType
	}
	return nil
}

func (ms *mediaService) removeOldUserPhoto(uploadsDir string, userID int) error {
	pattern := filepath.Join(uploadsDir, fmt.Sprintf("%d.*", userID))
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
		log.Println("Old photo deleted:", file)
	}
	return nil
}
