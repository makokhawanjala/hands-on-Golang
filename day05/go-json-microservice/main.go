package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ID      int     `json:"id,omitempty"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	InStock bool    `json:"in_stock"`
}

// POST /products → receive JSON and save to DB
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO products (name, price, in_stock) VALUES ($1, $2, $3) RETURNING id`
	err = db.QueryRow(query, p.Name, p.Price, p.InStock).Scan(&p.ID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// GET /products → fetch all products as JSON
func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, name, price, in_stock FROM products`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Price, &p.InStock)
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func main() {
	initDB()

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			createProduct(w, r)
		case "GET":
			getProducts(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
