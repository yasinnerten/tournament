package model

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type TournamentStatus string

const (
	Planned  TournamentStatus = "planned"
	Ongoing  TournamentStatus = "ongoing"
	Finished TournamentStatus = "finished"
)

type Tournament struct {
	ID     uint             `gorm:"primaryKey"`
	Name   string           `json:"name" validate:"required"`
	Status TournamentStatus `json:"status" validate:"required,default=planned"`
	Prize  int              `json:"prize" validate:"required"`
	Users  []User           `gorm:"many2many:tournament_users"`
}

func (t *Tournament) Validate() error {
	validate := validator.New()
	return validate.Struct(t)
}

// boş userlı ve planned turnuva oluşturmak için
func (t *Tournament) BeforeCreate(tx *gorm.DB) (err error) {
	t.Status = Planned
	t.Users = []User{}
	return
}
