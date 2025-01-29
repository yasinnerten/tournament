package crud

import (
	"tournament-app/internal/db"
	"tournament-app/model"
)

func CreateTournament(tournament *model.Tournament) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(tournament).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func GetOngoingTournaments() ([]model.Tournament, error) {
	var tournaments []model.Tournament
	if err := db.DB.Where("status = ?", "ongoing").Find(&tournaments).Error; err != nil {
		return nil, err
	}
	return tournaments, nil
}

func GetTournamentByID(id uint) (*model.Tournament, error) {
	var tournament model.Tournament
	if err := db.DB.Preload("Users").First(&tournament, id).Error; err != nil {
		return nil, err
	}
	return &tournament, nil
}

func UpdateTournament(tournament *model.Tournament) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(tournament).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func DeleteTournament(id uint) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&model.Tournament{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func GetAllTournaments() ([]model.Tournament, error) {
	var tournaments []model.Tournament
	if err := db.DB.Preload("Users").Find(&tournaments).Error; err != nil {
		return nil, err
	}
	return tournaments, nil
}

func GetTournamentByKey(key string) (*model.Tournament, error) {
	var tournament model.Tournament
	if err := db.DB.Where("key = ?", key).First(&tournament).Error; err != nil {
		return nil, err
	}
	return &tournament, nil
}

func UpdateLeaderboardEntry(entry *model.Leaderboard) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(entry).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func SaveLeaderboard(tournamentID uint, leaderboard []model.Leaderboard) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, entry := range leaderboard {
		entry.TournamentID = tournamentID
		if err := tx.Create(&entry).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func GetLeaderboardByTournamentID(tournamentID uint) ([]model.Leaderboard, error) {
	var leaderboard []model.Leaderboard
	if err := db.DB.Where("tournament_id = ?", tournamentID).Find(&leaderboard).Error; err != nil {
		return nil, err
	}
	return leaderboard, nil
}

func GetLeaderboardByUserID(userID uint) ([]model.Leaderboard, error) {
	var leaderboard []model.Leaderboard
	if err := db.DB.Where("user_id = ?", userID).Find(&leaderboard).Error; err != nil {
		return nil, err
	}
	return leaderboard, nil
}

// Redis-related functions will be retrieved from internal/db/redis.go
// CreateLeaderboardEntry creates a leaderboard entry in Redis
func CreateLeaderboardEntry(entry *model.Leaderboard) error {
	return db.CreateLeaderboardEntry(entry)
}

// GetLeaderboard retrieves the leaderboard from Redis
func GetLeaderboard(start, stop int64) ([]model.Leaderboard, error) {
	return db.GetLeaderboard(start, stop)
}

// UpdateLeaderboard updates the leaderboard in Redis
func UpdateLeaderboard(userID string, score float64) error {
	return db.UpdateLeaderboard(userID, score)
}

// RemoveLeaderboardFromRedis removes the leaderboard from Redis
func RemoveLeaderboardFromRedis(tournamentID uint) error {
	return db.RemoveLeaderboardFromRedis(tournamentID)
}
