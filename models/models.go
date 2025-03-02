package models

type User struct {
	ID       int     `json:"id,omitempty" gorm:"primaryKey"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8,max=32"`
	IsAdmin  bool    `json:"is_admin,omitempty" gorm:"default:false"`
	Profile  Profile `json:"profile,omitempty" gorm:"foreignKey:UserID" validate:"-"`
	Jobs     []Job   `json:"jobs,omitempty" gorm:"foreignKey:UserID" validate:"-"`
}

type Profile struct {
	ID          int    `json:"id,omitempty" gorm:"primaryKey"`
	UserID      int    `json:"user_id,omitempty"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=32"`
	LastName    string `json:"last_name" validate:"required,min=2,max=32"`
	PhoneNumber string `json:"phone_number" validate:"required,min=9,max=15"`
	PhotoPath   string `json:"-"`
	Header      string `json:"header" validate:"required,min=10,max=200"`
}

type Job struct {
	ID          int    `json:"id,omitempty" gorm:"primaryKey"`
	UserID      int    `json:"user_id,omitempty"`
	CompanyName string `json:"company_name" validate:"required,min=1,max=50"`
	Position    string `json:"position" validate:"required,min=2,max=50"`
	StartDate   string `json:"start_date" validate:"required"`
	EndDate     string `json:"end_date"`
	Description string `json:"description" validate:"required,min=10,max=500"`
}
