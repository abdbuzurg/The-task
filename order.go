package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type Order struct {
	ID        int  `json:"id"`
	UserID    int  `json:"userId"`
	ProductID int  `json:"productId"`
	Amount    uint `json:"amount"`
}

func InitialMigraitionForOrder() {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&Order{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to make initial migration for Order Table")
	}
}

// POST /order/
func createOrder(w http.ResponseWriter, r *http.Request) {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		RespondWithJson(w, "Invalid Body Fields", false)
		return
	}
	if order.Amount == 0 || order.ProductID == 0 || order.UserID == 0 {
		RespondWithJson(w, "Invalid Body Values", false)
		return
	}

	var user User
	if result := db.Select("balance").Find(&user, order.UserID); result.Error != nil {
		fmt.Println(result.Error.Error())
		RespondWithJson(w, "Invalid User ID", false)
		return
	}

	var productDetail ProductDetail
	if result := db.Select("amount_left", "price").Find(&productDetail, order.ProductID); result.Error != nil {
		RespondWithJson(w, "Invalid Product ID", false)
		return
	}

	if order.Amount > productDetail.AmountLeft {
		RespondWithJson(w, "We don't have the amount you want", false)
		return
	}

	if productDetail.Price*float32(order.Amount) > user.Balance {
		RespondWithJson(w, "You don't have enough balance", false)
		return
	}

	user.Balance = user.Balance - float32(order.Amount)*productDetail.Price
	productDetail.AmountLeft = productDetail.AmountLeft - order.Amount

	result := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ProductDetail{}).Where("id = ?", order.ProductID).Updates(productDetail).Error; err != nil {
			RespondWithJson(w, "Internal server error. Failed to update ProductDetail", false)
			return err
		}
		if err := tx.Model(&User{}).Where("id = ?", order.UserID).Updates(user).Error; err != nil {
			RespondWithJson(w, "Internal server error. Failed to update user balance", false)
			return err
		}
		if err := tx.Create(&order).Error; err != nil {
			RespondWithJson(w, "Internal server error. Failed to create new order", false)
			return err
		}

		return nil
	})

	if result != nil {
		RespondWithJson(w, "Failed to complete transaction", false)
		return
	}

	RespondWithJson(w, "New order has been made", true)
}
