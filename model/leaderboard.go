package model

import (
	"github.com/go-playground/validator/v10"
)

type LeaderboardStatus string

const (
	Active  LeaderboardStatus = "active"
	Passive LeaderboardStatus = "passive"
)

type Leaderboard struct {
	ID           uint              `gorm:"primaryKey"`
	UserID       uint              `json:"user_id" validate:"required"`
	TournamentID uint              `json:"tournament_id" validate:"required"`
	Score        float64           `json:"score" validate:"gte=0"`
	Status       LeaderboardStatus `json:"status" validate:"required"`
}

func (l *Leaderboard) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}
