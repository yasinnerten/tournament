package validation

import (
	"errors"

	"tournament-app/model"
)

func ValidateLeaderboard(leaderboard *model.Leaderboard) error {
	if leaderboard.UserID == 0 {
		return errors.New("leaderboard user_id cannot be empty")
	}
	if leaderboard.TournamentID == 0 {
		return errors.New("leaderboard tournament_id cannot be empty")
	}
	if leaderboard.Score < 0 {
		return errors.New("leaderboard score cannot be negative")
	}
	if leaderboard.Status != model.Active && leaderboard.Status != model.Passive {
		return errors.New("leaderboard status must be either 'active' or 'passive'")
	}

	return nil
}
