package validation

import (
	"errors"

	"tournament-app/model"
)

func ValidateTournament(tournament *model.Tournament) error {
	if tournament.Name == "" {
		return errors.New("tournament name cannot be empty")
	}
	if tournament.Prize < 0 {
		return errors.New("tournament prize cannot be negative")
	}
	if tournament.Status != model.Planned && tournament.Status != model.Ongoing && tournament.Status != model.Finished {
		return errors.New("tournament status must be either 'planned', 'ongoing', or 'finished'")
	}

	return nil
}
