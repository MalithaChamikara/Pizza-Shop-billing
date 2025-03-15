package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"piza_shop_billing/backend/models"
)

//function to handle POST requests to link pizza type and topping
func LinkPizzaTopping(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//extract the pizza type id from the request
		vars := mux.Vars(r)
		pizzaTypeId := vars["pizza_type_id"]


		//decode the request body to get the topping id
		var requestBody struct {
			ToppingId string `json:"topping_id"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Check if the pizza_type_id exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pizza_types WHERE pizza_type_id = ?)", pizzaTypeId).Scan(&exists)
		if err != nil || !exists {
		 http.Error(w, "Pizza type not found", http.StatusBadRequest)
			return
		}

		// Check if the topping_id exists
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM toppings WHERE topping_id = ?)", requestBody.ToppingId).Scan(&exists)
		if err != nil || !exists {
		 http.Error(w, "Topping not found", http.StatusBadRequest)
		 return
		}

		//query to insert pizza type and topping into pizza_toppings table
		query := `INSERT INTO pizza_toppings (pizza_type_id, topping_id) VALUES (?, ?)`

		if _, err := db.Exec(query,pizzaTypeId,requestBody.ToppingId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the response object
		pizzaTopping := models.PizzaTopping{
        	PizzaTypeId: pizzaTypeId,
            ToppingId:   requestBody.ToppingId,
        }

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(pizzaTopping)
	}
}

//function to get topping names related to specific pizzatype
// Function to get topping names related to a specific pizza type
func GetToppingsByPizzaType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the pizza type ID from the request
		vars := mux.Vars(r)
		pizzaTypeId := vars["pizza_type_id"]

		// Check if the pizza_type_id exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pizza_types WHERE pizza_type_id = ?)", pizzaTypeId).Scan(&exists)
		if err != nil || !exists {
			http.Error(w, "Pizza type not found", http.StatusBadRequest)
			return
		}

		// Query to get the topping names for the given pizza type ID
		query := `
			SELECT t.name 
			FROM toppings t
			INNER JOIN pizza_toppings pt ON t.topping_id = pt.topping_id
			WHERE pt.pizza_type_id = ?
		`

		// Execute the query
		rows, err := db.Query(query, pizzaTypeId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Collect the topping names
		var toppings []string
		for rows.Next() {
			var toppingName string
			if err := rows.Scan(&toppingName); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			toppings = append(toppings, toppingName)
		}

		// Check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return the topping names as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toppings)
	}
}
