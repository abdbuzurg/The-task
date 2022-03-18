package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	db.Model(&ProductDetail{}).Where("id = ?", order.ProductID).Updates(productDetail)
	db.Model(&User{}).Where("id = ?", order.UserID).Updates(user)
	db.Create(&order)

	RespondWithJson(w, "New order has been made", true)
}
