package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Starting server...")

	//Table Initialization
	InitialMigraitionForCategory()
	InitialMigraitionForProductDetail()
	InitialMigraitionForProduct()
	InitialMigraitionForUser()
	InitialMigraitionForOrder()

	router := mux.NewRouter()
	router.StrictSlash(true)

	// CRUD for categories
	router.HandleFunc("/category/", GetAllCategories).Methods("GET")
	router.HandleFunc("/category/{id}", GetCategory).Methods("GET")
	router.HandleFunc("/category/", CreateCategory).Methods("POST")
	router.HandleFunc("/category/{id}", UpdateCategory).Methods("PATCH", "PUT")
	router.HandleFunc("/category/{id}", DeleteCategory).Methods("DELETE")

	// CRUD for products
	router.HandleFunc("/product/", GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", GetProduct).Methods("GET")
	router.HandleFunc("/product/", CreateProducts).Methods("POST")
	router.HandleFunc("/product/{id}", UpdateProducts).Methods("PUT", "PATCH")
	router.HandleFunc("/product/{id}", DeleteProducts).Methods("DELETE")

	// Crate for order
	router.HandleFunc("/order/", createOrder).Methods("POST")

	fmt.Println("Server successfully started at port 8080")
	err := http.ListenAndServe("127.0.0.1:8080", router)
	if err != nil {
		panic(err)
	}

}
