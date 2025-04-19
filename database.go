package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Quote struct {
	gorm.Model
	Content string
	Author  string
	Votes   []Vote
}

type Vote struct {
	gorm.Model
	PersonID int64
	QuoteID  uint
}

func loadDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("quotobot.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Quote{}, &Vote{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database loaded successfully")

	return db
}
