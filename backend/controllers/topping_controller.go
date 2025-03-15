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

//method to get all toppings and  returns a http.HandlerFunc
func GetToppings(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := " SELECT topping_id,name,price FROM toppings"
		results, err := db.Query(query)

		//check if there is an error and return it to the client
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//defer is used to close the result set after the function returns
		defer results.Close()

		//create a slice of toppings
		var toppings []models.Topping
		//iterate over the result set and append the toppings to the slice
		for results.Next() {
			var topping models.Topping
			//scan the result set into the topping struct and check for errors
			if err := results.Scan(&topping.ToppingId, &topping.Name,&topping.Price); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//append the topping to the slice
			toppings = append(toppings, topping)
		}

		//set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toppings)
	}
}

//method to create a topping
func CreateTopping(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//create a topping struct
		var topping models.Topping
		//decode the request body into the topping struct
		if err := json.NewDecoder(r.Body).Decode(&topping); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		topping.CreatedAt = time.Now()
		topping.UpdatedAt = time.Now()
		//create a query to insert the topping into the database
		query := "INSERT INTO toppings(topping_id,name,price,created_at,updated_at) VALUES(?,?,?,?,?)"
		//execute the query and check for errors
		if _, err := db.Exec(query, topping.ToppingId,topping.Name,topping.Price,topping.CreatedAt,topping.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		//encode the topping struct into json and write it to the response writer
		json.NewEncoder(w).Encode(topping)
	}
}

//method to update a topping
func UpdateTopping(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract topping_id from URL path
        vars := mux.Vars(r)
        toppingId := vars["topping_id"]

        // Fetch the existing record from the database
        var existingTopping models.Topping
        query := "SELECT topping_id,name,price FROM toppings WHERE topping_id = ?"
        if err := db.QueryRow(query, toppingId).Scan(
            &existingTopping.ToppingId,
            &existingTopping.Name,
            &existingTopping.Price,
        ); err != nil {
            log.Printf("Error fetching existing Topping: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }


        // Decode the request body into a new topping struct
        var updatedTopping models.Topping
        if err := json.NewDecoder(r.Body).Decode(&updatedTopping); err != nil {
            log.Printf("Error decoding request body: %v", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Merge the changes
        if updatedTopping.Name != "" {
            existingTopping.Name = updatedTopping.Name
        }
        if updatedTopping.Price != 0  {
            existingTopping.Price = updatedTopping.Price
        }
        
        existingTopping.UpdatedAt = time.Now()

        // Build the update query dynamically
        updateQuery := "UPDATE toppings SET name=?,price=?,updated_at=? WHERE topping_id=?"
        args := []interface{}{
            existingTopping.Name,
			existingTopping.Price,
			existingTopping.UpdatedAt,
            toppingId,
        }

        // Execute the query and check for errors
        if _, err := db.Exec(updateQuery, args...); err != nil {
            log.Printf("Error executing update query: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Set the response header to application/json
        w.Header().Set("Content-Type", "application/json")
        // Encode the updated topping struct into json and write it to the response writer
        if err := json.NewEncoder(w).Encode(existingTopping); err != nil {
            log.Printf("Error encoding response: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}
//method to delete a topping
func DeleteTopping(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract topping_id from the request URL
		vars := mux.Vars(r)
		toppingId := vars["topping_id"]

		//query to delete topping from associated tables
		relatedQuery:= "DELETE FROM pizza_toppings WHERE topping_id = ?"
		if _, err := db.Exec(relatedQuery, toppingId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// create a query to delete the topping from the database
		query := "DELETE FROM toppings WHERE topping_id = ?"
		// execute the query and check for errors
		if _, err := db.Exec(query, toppingId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		// write a success message to the response writer
		json.NewEncoder(w).Encode(map[string]string{"message": "Topping deleted successfully"})
	}
}