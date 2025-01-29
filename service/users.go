package service

import (
	"tournament-app/internal/crud"
	"tournament-app/model"
	"tournament-app/validation"
)

// CreateUser validates and creates a new user
func CreateUser(user *model.User) error {
	if err := validation.ValidateUser(user); err != nil {
		return err
	}
	return crud.CreateUser(user)
}

// UpdateUser validates and updates an existing user
func UpdateUser(user *model.User) error {
	if err := validation.ValidateUser(user); err != nil {
		return err
	}
	return crud.UpdateUser(user)
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id uint) (*model.User, error) {
	return crud.GetUserByID(id)
}

// GetUsers retrieves all users
func GetUsers() ([]model.User, error) {
	return crud.GetUsers()
}

func PerformHealthCheck() (string, error) {
	// Check PostgreSQL connection
	if err := crud.PingPostgres(); err != nil {
		return "PostgreSQL is not healthy", err
	}

	// Check Redis connection
	if err := crud.PingRedis(); err != nil {
		return "Redis is not healthy", err
	}

	return "Service is healthy", nil
}

func DeleteUser(id uint) error {
	return crud.DeleteUser(id)
}

func ClearDatabase() error {
	return crud.ClearDatabase()
}
