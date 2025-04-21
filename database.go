package main

import (
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
	PersonID int64 `gorm:"uniqueIndex:udx_person_quote,WHERE:deleted_at IS NULL"`
	QuoteID  uint  `gorm:"uniqueIndex:udx_person_quote,WHERE:deleted_at IS NULL"`
}

func loadDatabase(logger *Logger) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("quotobot.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		logger.Error.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Quote{}, &Vote{}); err != nil {
		logger.Error.Fatalf("Failed to migrate database: %v", err)
	}

	logger.Info.Println("Database loaded successfully")

	return db
}
