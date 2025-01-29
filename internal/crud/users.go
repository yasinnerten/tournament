package crud

import (
	"log"
	"tournament-app/internal/db"
	"tournament-app/model"
)

// Testing connection to the database
func PingPostgres() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

func PingRedis() error {
	return db.PingRedis()
}

// User related functions
func CreateUser(user *model.User) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		log.Printf("Error creating user: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	log.Printf("User created successfully: %v", user)
	return nil
}

func UpdateUser(user *model.User) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		log.Printf("Error updating user: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers() ([]model.User, error) {
	var users []model.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(id uint) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&model.User{}, id).Error; err != nil {
		tx.Rollback()
		log.Printf("Error deleting user: %v", err)
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// General function to clear the database
func ClearDatabase() error {
	if err := db.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}

	if err := db.DB.Exec("TRUNCATE TABLE tournaments RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}

	if err := db.DB.Exec("TRUNCATE TABLE leaderboards RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}

	return nil
}
