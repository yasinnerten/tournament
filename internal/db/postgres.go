package db

import (
	"log"
	"tournament-app/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitPostgres initializes the PostgreSQL database
func InitPostgres(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to Postgres: %v", err)
		return err
	}
	log.Println("Successfully connected to Postgres")

	// migrateSchema fonksiyonu burda kalsın lazım olur
	// err = migrateSchema()
	// if err != nil {
	// 	return err

	err = DB.AutoMigrate(
		&model.User{},
		&model.Tournament{},
		&model.Leaderboard{},
	)

	if err != nil {
		log.Printf("Failed to migrate schema: %v", err)
		return err
	}

	return nil
}
