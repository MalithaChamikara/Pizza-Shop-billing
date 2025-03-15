package controllers

import (
	"database/sql"
	"encoding/json"
	"time"
	"log"
	"net/http"
	"piza_shop_billing/backend/models"

	"github.com/gorilla/mux"
)

//method to get all pizza types returns a http.HandlerFunc
func GetPizzaTypes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT pizza_type_id,name,size,base_price,description FROM pizza_types"
		results, err := db.Query(query)

		//check if there is an error and return it to the client
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//defer is used to close the result set after the function returns
		defer results.Close()

		//create a slice of pizza types
		var pizzaTypes []models.PizzaType
		//iterate over the result set and append the pizza types to the slice
		for results.Next() {
			var pizzaType models.PizzaType
			//scan the result set into the pizza type struct and check for errors
			if err := results.Scan(&pizzaType.PizzaTypeId, &pizzaType.Name, &pizzaType.Size, &pizzaType.BasePrice, &pizzaType.Description); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//append the pizza type to the slice
			pizzaTypes = append(pizzaTypes, pizzaType)
		}

		//set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pizzaTypes)
	}
}

//method to create a pizza type
func CreatePizzaType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//create a pizza type struct
		var pizzaType models.PizzaType
		//decode the request body into the pizza type struct
		if err := json.NewDecoder(r.Body).Decode(&pizzaType); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pizzaType.CreatedAt = time.Now()
		pizzaType.UpdatedAt = time.Now()
		//create a query to insert the pizza type into the database
		query := "INSERT INTO pizza_types(pizza_type_id,name,size,base_price,description,created_at,updated_at) VALUES(?,?,?,?,?,?,?)"
		//execute the query and check for errors
		if _, err := db.Exec(query, pizzaType.PizzaTypeId,pizzaType.Name, pizzaType.Size, pizzaType.BasePrice, pizzaType.Description,pizzaType.CreatedAt,pizzaType.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		//encode the pizza type struct into json and write it to the response writer
		json.NewEncoder(w).Encode(pizzaType)
	}
}

//method to update a pizza type
func UpdatePizzaType(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract pizza_type_id from URL path
        vars := mux.Vars(r)
        pizzaTypeId := vars["pizza_type_id"]

        // Fetch the existing record from the database
        var existingPizzaType models.PizzaType
        query := "SELECT pizza_type_id, name, size, base_price, description FROM pizza_types WHERE pizza_type_id = ?"
        if err := db.QueryRow(query, pizzaTypeId).Scan(
            &existingPizzaType.PizzaTypeId,
            &existingPizzaType.Name,
            &existingPizzaType.Size,
            &existingPizzaType.BasePrice,
            &existingPizzaType.Description,
        ); err != nil {
            log.Printf("Error fetching existing pizza type: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }


        // Decode the request body into a new pizza type struct
        var updatedPizzaType models.PizzaType
        if err := json.NewDecoder(r.Body).Decode(&updatedPizzaType); err != nil {
            log.Printf("Error decoding request body: %v", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Merge the changes
        if updatedPizzaType.Name != "" {
            existingPizzaType.Name = updatedPizzaType.Name
        }
        if updatedPizzaType.Size != "" {
            existingPizzaType.Size = updatedPizzaType.Size
        }
        if updatedPizzaType.BasePrice != 0 {
            existingPizzaType.BasePrice = updatedPizzaType.BasePrice
        }
        if updatedPizzaType.Description != "" {
            existingPizzaType.Description = updatedPizzaType.Description
        }
        existingPizzaType.UpdatedAt = time.Now()

        // Build the update query dynamically
        updateQuery := "UPDATE pizza_types SET name=?, size=?, base_price=?, description=?, updated_at=? WHERE pizza_type_id=?"
        args := []interface{}{
            existingPizzaType.Name,
            existingPizzaType.Size,
            existingPizzaType.BasePrice,
            existingPizzaType.Description,
            existingPizzaType.UpdatedAt,
            pizzaTypeId,
        }

        // Execute the query and check for errors
        if _, err := db.Exec(updateQuery, args...); err != nil {
            log.Printf("Error executing update query: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Set the response header to application/json
        w.Header().Set("Content-Type", "application/json")
        // Encode the updated pizza type struct into json and write it to the response writer
        if err := json.NewEncoder(w).Encode(existingPizzaType); err != nil {
            log.Printf("Error encoding response: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}
//method to delete a pizza type
func DeletePizzaType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract pizza_type_id from the request URL
		vars := mux.Vars(r)
		pizzaTypeId := vars["pizza_type_id"]

		//query to delete topping from associated tables
		relatedQuery:= "DELETE FROM pizza_toppings WHERE pizza_type_id = ?"
		if _, err := db.Exec(relatedQuery, pizzaTypeId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}
		// create a query to delete the pizza type from the database
		query := "DELETE FROM pizza_types WHERE pizza_type_id = ?"
		// execute the query and check for errors
		if _, err := db.Exec(query, pizzaTypeId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		// write a success message to the response writer
		json.NewEncoder(w).Encode(map[string]string{"message": "Pizza type deleted successfully"})
	}
}