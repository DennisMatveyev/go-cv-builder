package user

import (
	"errors"
	"log"

	"cv-builder/common"
	"cv-builder/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	GetProfile(userID int) (*models.Profile, error)
	GetJobs(userID int) ([]*models.Job, error)
	CreateProfile(userID int, profile *models.Profile) (int, error)
	CreateJob(userID int, job *models.Job) (int, error)
	UpdateProfile(profile *models.Profile) error
	UpdateJob(job *models.Job) error
	DeleteJob(jobID int) error
	GetUserData(userID int) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SaveUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var userDB models.User
	err := r.db.Where("email = ?", email).First(&userDB).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("User with email %s not found", email)
		return nil, nil
	}
	if err != nil {
		log.Printf("Database error when finding user by email %s: %v", email, err)
		return nil, err
	}
	return &userDB, nil
}

func (r *userRepository) GetProfile(userID int) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.Model(&models.Profile{}).Where("user_id = ?", userID).First(&profile).Error
	if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return &profile, nil
	} else {
		log.Printf("Database error when getting profile by userID %d: %v", userID, err)
		return nil, common.ErrDatabase
	}
}

func (r *userRepository) GetJobs(userID int) ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.Model(&models.Job{}).Where("user_id = ?", userID).Find(&jobs).Error
	if err != nil {
		log.Printf("Database error when finding jobs for user ID %d: %v", userID, err)
		return nil, common.ErrDatabase
	}
	return jobs, nil
}

func (r *userRepository) CreateProfile(userID int, profile *models.Profile) (int, error) {
	profileFound := new(models.Profile)
	r.db.Model(&models.Profile{}).Where("user_id = ?", userID).First(&profileFound)
	if profileFound.ID != 0 {
		log.Printf("Profile already exists for user ID %d", userID)
		return 0, ErrProfileExists
	}
	profile.UserID = userID
	err := r.db.Create(profile).Error
	if err != nil {
		log.Printf("Database error when creating profile for user ID %d: %v", userID, err)
		return 0, common.ErrDatabase
	}
	return profile.ID, nil
}

func (r *userRepository) CreateJob(userID int, job *models.Job) (int, error) {
	job.UserID = userID
	err := r.db.Create(job).Error
	if err != nil {
		log.Printf("Database error when creating job for user ID %d: %v", userID, err)
		return 0, common.ErrDatabase
	}
	return job.ID, nil
}

func (r *userRepository) UpdateProfile(profile *models.Profile) error {
	err := r.db.Save(profile).Error
	if err != nil {
		log.Printf("Database error when updating profile with ID %d: %v", profile.ID, err)
		return common.ErrDatabase
	}
	return nil
}

func (r *userRepository) UpdateJob(job *models.Job) error {
	err := r.db.Save(job).Error
	if err != nil {
		log.Printf("Database error when updating job with ID %d: %v", job.ID, err)
		return common.ErrDatabase
	}
	return nil
}

func (r *userRepository) DeleteJob(jobID int) error {
	err := r.db.Delete(&models.Job{}, jobID).Error
	if err != nil {
		log.Printf("Database error when deleting job with ID %d: %v", jobID, err)
		return common.ErrDatabase
	}
	return nil
}

func (r *userRepository) GetUserData(userID int) (*models.User, error) {
	user := new(models.User)
	if err := r.db.Preload("Profile").Preload("Jobs").First(user, userID).Error; err != nil {
		return nil, common.ErrDatabase
	}
	return user, nil
}
