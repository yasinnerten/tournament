package service

import (
	"fmt"
	"tournament-app/internal/crud"
	"tournament-app/model"
	"tournament-app/validation"
)

func CreateTournament(tournament *model.Tournament) error {
	tournament.Status = model.Planned
	if err := crud.CreateTournament(tournament); err != nil {
		return err
	}

	// Create an empty leaderboard for the each tournament created
	leaderboard := model.Leaderboard{
		TournamentID: tournament.ID,
	}
	if err := CreateLeaderboardEntry(&leaderboard); err != nil {
		return err
	}

	return nil
}

func UpdateTournament(tournament *model.Tournament) error {
	if err := validation.ValidateTournament(tournament); err != nil {
		return err
	}
	return crud.UpdateTournament(tournament)
}

func DeleteTournament(id uint) error {
	return crud.DeleteTournament(id)
}

func GetTournamentByID(id uint) (*model.Tournament, error) {
	return crud.GetTournamentByID(id)
}

func GetAllTournaments() ([]model.Tournament, error) {
	return crud.GetAllTournaments()
}

func GetOngoingTournaments() ([]model.Tournament, error) {
	return crud.GetOngoingTournaments()
}

func EndTournament(tournamentID uint) error {
	tournament, err := crud.GetTournamentByID(tournamentID)
	if err != nil {
		return err
	}

	//10 kişiden fazla katılım olursa ya da turnuva elle bitirilirse
	if len(tournament.Users) >= 10 || tournament.Status == model.Finished {
		tournament.Status = model.Finished
		if err := crud.UpdateTournament(tournament); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("tournament cannot be ended")
}

// JoinTournament allows a user to join a tournament
func JoinTournament(tournamentID, userID uint) error {
	tournament, err := crud.GetTournamentByID(tournamentID)
	if err != nil {
		return err
	}

	// Check if the tournament is already finished
	if tournament.Status == model.Finished {
		return fmt.Errorf("cannot join a finished tournament")
	}

	user, err := crud.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Decrease user's money by 50
	if user.Money < 50 {
		return fmt.Errorf("user does not have enough money to join the tournament")
	}
	user.Money -= 50
	if err := crud.UpdateUser(user); err != nil {
		return err
	}

	tournament.Users = append(tournament.Users, *user)
	if len(tournament.Users) >= 10 {
		tournament.Status = model.Finished // Finish the tournament if user count exceeds 10
		if err := SetLeaderboard("leaderboard:"+tournament.Name, tournament.Users); err != nil {
			return err
		}
		if err := FinalizeTournament("leaderboard:" + tournament.Name); err != nil {
			return err
		}
	}
	if err := crud.UpdateTournament(tournament); err != nil {
		return err
	}

	// Update leaderboard in Redis when someone joins the tournament
	score := calculateScore(user)
	if err := crud.UpdateLeaderboard(fmt.Sprintf("%d", user.ID), score); err != nil {
		return err
	}

	return nil
}

// Status active olan leaderboardları görmek için
func GetActiveLeaderboard(start, stop int64) ([]model.Leaderboard, error) {
	leaderboard, err := crud.GetLeaderboard(start, stop)
	if err != nil {
		return nil, err
	}

	var activeLeaderboard []model.Leaderboard
	for _, entry := range leaderboard {
		if entry.Status == model.Active {
			activeLeaderboard = append(activeLeaderboard, entry)
		}
	}
	return activeLeaderboard, nil
}

// bir kişinin katıldığı active leaderboardları görmek için
func GetActiveLeaderboardByUserID(userID uint) ([]model.Leaderboard, error) {
	leaderboard, err := crud.GetLeaderboardByUserID(userID)
	if err != nil {
		return nil, err
	}

	var activeLeaderboard []model.Leaderboard
	for _, entry := range leaderboard {
		if entry.Status == model.Active {
			activeLeaderboard = append(activeLeaderboard, entry)
		}
	}
	return activeLeaderboard, nil
}

// Turnuvaya ait active leaderboardu görmek için
func GetActiveLeaderboardByTournamentID(tournamentID uint) ([]model.Leaderboard, error) {
	leaderboard, err := crud.GetLeaderboardByTournamentID(tournamentID)
	if err != nil {
		return nil, err
	}

	var activeLeaderboard []model.Leaderboard
	for _, entry := range leaderboard {
		if entry.Status == model.Active {
			activeLeaderboard = append(activeLeaderboard, entry)
		}
	}
	return activeLeaderboard, nil
}

// Turnuvaya ait passive-bitmiş leaderboardu görmek için
func GetFinishedLeaderboardByTournamentID(tournamentID uint) ([]model.Leaderboard, error) {
	leaderboard, err := crud.GetLeaderboardByTournamentID(tournamentID)
	if err != nil {
		return nil, err
	}

	var finishedLeaderboard []model.Leaderboard
	for _, entry := range leaderboard {
		if entry.Status == model.Passive {
			finishedLeaderboard = append(finishedLeaderboard, entry)
		}
	}
	return finishedLeaderboard, nil
}

func FinalizeTournament(key string) error {

	tournament, err := crud.GetTournamentByKey(key)
	if err != nil {
		return fmt.Errorf("failed to retrieve tournament: %v", err)
	}

	// Check if the tournament is active
	if tournament.Status != model.Ongoing {
		return fmt.Errorf("tournament is not active (status=ongoing)")
	}

	// Retrieve the leaderboard for the tournament
	leaderboard, err := GetActiveLeaderboardByTournamentID(tournament.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve leaderboard: %v", err)
	}

	// Distribute prizes based on the leaderboard standings
	for i, entry := range leaderboard {
		user, err := crud.GetUserByID(entry.UserID)
		if err != nil {
			return fmt.Errorf("failed to retrieve user: %v", err)
		}

		// Calculate prize based on position
		prize := calculatePrize(tournament.Prize, i+1)
		user.Money += prize

		// Update user
		if err := crud.UpdateUser(user); err != nil {
			return fmt.Errorf("failed to update user: %v", err)
		}
	}

	// Save the leaderboard to PostgreSQL
	if err := crud.SaveLeaderboard(tournament.ID, leaderboard); err != nil {
		return fmt.Errorf("failed to save leaderboard: %v", err)
	}

	// Update leaderboard status to passive
	for _, entry := range leaderboard {
		entry.Status = model.Passive
		if err := crud.UpdateLeaderboardEntry(&entry); err != nil {
			return fmt.Errorf("failed to update leaderboard entry: %v", err)
		}
	}

	// Remove the leaderboard from Redis
	if err := crud.RemoveLeaderboardFromRedis(tournament.ID); err != nil {
		return fmt.Errorf("failed to remove leaderboard from Redis: %v", err)
	}

	// Update the tournament status to closed
	tournament.Status = model.Finished
	if err := crud.UpdateTournament(tournament); err != nil { // pointer used for tournament
		return fmt.Errorf("failed to update tournament: %v", err)
	}

	return nil
}

// calculatePrize calculates the prize based on the total prize pool and the position
func calculatePrize(totalPrize, position int) int {
	switch position {
	case 1:
		return totalPrize / 2 // 50% for 1st place
	case 2:
		return totalPrize / 4 // 25% for 2nd place
	case 3:
		return totalPrize / 8 // 12.5% for 3rd place
	default:
		return totalPrize / 16 // 6.25% for other places
	}
}

func SetLeaderboard(key string, users []model.User) error {
	for _, user := range users {
		score := calculateScore(&user)
		if err := crud.CreateLeaderboardEntry(&model.Leaderboard{
			UserID: user.ID,
			Score:  score,
			Status: model.Active,
		}); err != nil {
			return err
		}
	}
	return nil
}

func LevelUpUser(userID uint) error {
	user, err := crud.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Calculate the cost to level up
	cost := 100 + (user.Level * 50)

	if user.Money < cost {
		return fmt.Errorf("insufficient funds")
	}

	// Deduct the cost and increase the user's level
	user.Money -= cost
	user.Level += 1

	// Recalculate the user's score
	user.Score = calculateScore(user)

	// Update the user's data in PostgreSQL
	if err := crud.UpdateUser(user); err != nil {
		return err
	}

	// Update the leaderboard in Redis using the CRUD function
	if err := crud.UpdateLeaderboard(fmt.Sprintf("%d", user.ID), user.Score); err != nil {
		return err
	}

	return nil
}

// calculateScore calculates the score for a user based on their level and other factors
func calculateScore(user *model.User) float64 {
	return float64(user.Level*100 + user.Money)
}

func CreateLeaderboardEntry(entry *model.Leaderboard) error {
	// If UserID is not provided, skip user-related operations
	if entry.UserID == 0 {
		return crud.CreateLeaderboardEntry(entry)
	}

	user, err := crud.GetUserByID(entry.UserID)
	if err != nil {
		return err
	}

	// Calculate score based on user's level and money every time a new leaderboard is created
	entry.Score = calculateScore(user)

	return crud.CreateLeaderboardEntry(entry)
}
