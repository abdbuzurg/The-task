package main

import (
	"fmt"

	"gorm.io/gorm"
)

type ProductDetail struct {
	gorm.Model
	ID         uint    `json:"id"`
	Price      float32 `json:"price"`
	AmountLeft uint    `json:"amountLeft"`

	Product Product `gorm:"constraint:OnDelete:CASCADE"`
}

func InitialMigraitionForProductDetail() {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&ProductDetail{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to make initial migration for ProductDetail Table")
	}
}
