package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID              uint
	CategoryID      uint   `json:"categoryId"`
	ProductDetailID uint   `json:"productDetailId" gorm:"constraint:OnDelete:CASCADE"`
	Name            string `json:"name"`

	Order []Order
}

type ProductRequestBodyFormat struct {
	ID         uint    `json:"ID"`
	CategoryID uint    `json:"categoryId"`
	AmountLeft uint    `json:"amountLeft"`
	Price      float32 `json:"price"`
	Name       string  `json:"name"`
}

func InitialMigraitionForProduct() {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&Product{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to make initial migration for Product Table")
	}
}

// /product/ GET
func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	var products []Product
	result := db.Find(&products)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		RespondWithJson(w, "Internal Server Error", false)
		return
	}

	RespondWithJson(w, products, true)
}

// /product/{id} GET
func GetProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	id := mux.Vars(r)["id"]
	if _, err := strconv.Atoi(id); err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid URL parameters", false)
		return
	}

	var product Product
	result := db.First(&product, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		RespondWithJson(w, "No such entry exist", false)
		return
	}

	RespondWithJson(w, product, true)
}

// /product/ POST
func CreateProducts(w http.ResponseWriter, r *http.Request) {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	var newProduct ProductRequestBodyFormat
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid Body parameters", false)
		return
	}
	if isValidRequestBody(newProduct) {
		RespondWithJson(w, "Invalid Body values", false)
		return
	}

	result := db.Create(&ProductDetail{
		AmountLeft: newProduct.AmountLeft,
		Price:      newProduct.Price,
		Product: Product{
			Name:       newProduct.Name,
			CategoryID: newProduct.CategoryID,
		},
	})
	if result.Error != nil {
		RespondWithJson(w, "Internal Server Error", false)
		return
	}

	RespondWithJson(w, "New entry has been added", true)
}

// /product/{id} PUT or PATCH
func UpdateProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	id := mux.Vars(r)["id"]
	if _, err := strconv.Atoi(id); err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid URL parameters", false)
		return
	}

	var newProduct ProductRequestBodyFormat
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid Body parameters", false)
		return
	}
	if isValidRequestBody(newProduct) {
		RespondWithJson(w, "Invalid Body values", false)
		return
	}

	result := db.Model(&ProductDetail{}).Where("id = ?", id).Updates(ProductDetail{
		AmountLeft: newProduct.AmountLeft,
		Price:      newProduct.Price,
	})
	if result.Error != nil {
		RespondWithJson(w, "Internal Server Error", false)
		return
	}
	result = db.Model(&Product{}).Where("id = ?", id).Updates(Product{
		Name:       newProduct.Name,
		CategoryID: newProduct.CategoryID,
	})
	if result.Error != nil {
		RespondWithJson(w, "Internal Server Error", false)
		return
	}

	RespondWithJson(w, newProduct, true)
}

// /product/{id} DELETE
func DeleteProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	id := mux.Vars(r)["id"]
	if _, err := strconv.Atoi(id); err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid URL parameters", false)
		return
	}

	result := db.Delete(&Product{}, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		RespondWithJson(w, "No such entry exist", false)
		return
	}

	RespondWithJson(w, "An entry has been deleted", true)
}

func isValidRequestBody(product ProductRequestBodyFormat) bool {
	return product.Name == "" || product.CategoryID == 0 || product.Price == 0 || product.AmountLeft == 0
}
