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

type Category struct {
	gorm.Model
	ID   uint
	Name string `gorm:"unique, not null" json:"name"`

	//relationships
	Product []Product
}

// Initial Migration Setup
func InitialMigraitionForCategory() {
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	err := db.AutoMigrate(&Category{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to make initial migration for Category Table")
	}
}

// GET /category/
func GetAllCategories(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	var categories []Category
	result := db.Find(&categories)

	if result.Error != nil {
		RespondWithJson(w, "Internal Server Error", false)
		return
	}

	RespondWithJson(w, categories, true)
}

// GET /category/{id}
func GetCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	id := mux.Vars(r)["id"]
	if _, err := strconv.Atoi(id); err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid URL parameters", false)
		return
	}

	var category Category
	result := db.First(&category, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		RespondWithJson(w, "No such entry exist", false)
		return
	}

	RespondWithJson(w, category, true)
}

// POST /category/
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	var categeory Category
	err := json.NewDecoder(r.Body).Decode(&categeory)
	if err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid Body Parameters", false)
		return
	}
	if categeory.Name == "" {
		RespondWithJson(w, "Invalid Body values", false)
		return
	}

	result := db.Create(&categeory)

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		RespondWithJson(w, "Internal Server Error", true)
	}

	RespondWithJson(w, categeory, true)
}

//  PUT or PATCH /category/{id}
func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	id := mux.Vars(r)["id"]
	if _, err := strconv.Atoi(id); err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid URL parameters", false)
		return
	}

	var newCategeory Category
	err := json.NewDecoder(r.Body).Decode(&newCategeory)
	if err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid Body Parameters", false)
		return
	}
	if newCategeory.Name == "" {
		RespondWithJson(w, "Invalid Body values", false)
		return
	}

	result := db.Model(&Category{}).Where("id = ?", id).Update("name", newCategeory.Name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		RespondWithJson(w, "No such entry exist", false)
		return
	}

	RespondWithJson(w, newCategeory, true)
}

// DELETE /category/{id}
func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path + " " + r.Method)
	db, sqlDB := GetDB()
	defer sqlDB.Close()

	id := mux.Vars(r)["id"]
	if _, err := strconv.Atoi(id); err != nil {
		fmt.Println(err.Error())
		RespondWithJson(w, "Invalid URL parameters", false)
		return
	}

	result := db.Delete(&Category{}, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		RespondWithJson(w, "No such entry exist", false)
		return
	}

	RespondWithJson(w, "An entry from Category tables has been deleted", true)
}
