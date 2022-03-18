package main

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Balance  float32 `json:"balance"`

	Order []Order
}

func InitialMigraitionForUser() {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to make initial migration for Category Table")
	}
}
