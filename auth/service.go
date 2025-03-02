package auth

import (
	"log"
	"time"

	"cv-builder/common"
	"cv-builder/models"
	"cv-builder/user"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo  user.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo user.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (as *AuthService) Register(user *models.User) error {
	userFound, err := as.userRepo.FindByEmail(user.Email)
	if userFound != nil {
		return ErrUserExists
	}
	if err != nil {
		return common.ErrDatabase
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if err := as.userRepo.SaveUser(user); err != nil {
		return ErrSaveUser
	}

	return nil
}

func (as *AuthService) Login(email, password string) (string, error) {
	userDB, err := as.userRepo.FindByEmail(email)

	if userDB != nil && PasswordValid(userDB.Password, password) {
		return as.generateToken(userDB.ID)
	} else if userDB == nil && err == nil {
		return "", ErrInvalidCredentials
	} else {
		return "", common.ErrDatabase
	}
}

func (as *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(as.jwtSecret))
	if err != nil {
		log.Println("Failed to generate token: ", err)
		return "", ErrGenerateToken
	}

	return token, nil
}
