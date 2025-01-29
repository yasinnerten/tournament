package model

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type User struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `json:"name" validate:"required"`
	Money int     `json:"money" validate:"required"`
	Level int     `json:"level" validate:"required"`
	Score float64 `json:"score" validate:"gte=0"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// BeforeCreate hook to calculate the score before saving the user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Score = calculateScore(u)
	return
}

// calculateScore calculates the score for a user based on their level and other factors
func calculateScore(user *User) float64 {
	return float64(user.Level*100 + user.Money)
}
