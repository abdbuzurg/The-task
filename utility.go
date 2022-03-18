package main

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB connection
func GetDB() (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	sqlDB, _ := db.DB()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to Database")
	}
	return db, sqlDB
}
