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

//method to get all beverages returns a http.HandlerFunc
func GetBeverages(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT beverage_id,name,price FROM beverages"
		results, err := db.Query(query)

		//check if there is an error and return it to the client
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//defer is used to close the result set after the function returns
		defer results.Close()

		//create a slice of beverages
		var beverages []models.Beverage
		//iterate over the result set and append the beverages to the slice
		for results.Next() {
			var beverage models.Beverage
			//scan the result set into the beverage struct and check for errors
			if err := results.Scan(&beverage.BeverageId, &beverage.Name, &beverage.Price); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//append the beverage to the slice
			beverages = append(beverages, beverage)
		}

		//set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(beverages)
	}
}

//method to create a beverage
func CreateBeverage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//create a beverage struct
		var beverage models.Beverage
		//decode the request body into the beverage struct
		if err := json.NewDecoder(r.Body).Decode(&beverage); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		beverage.CreatedAt = time.Now()
		beverage.UpdatedAt = time.Now()
		//create a query to insert the beverage into the database
		query := "INSERT INTO beverages(beverage_id,name,price,created_at,updated_at) VALUES(?,?,?,?,?)"
		//execute the query and check for errors
		if _, err := db.Exec(query, beverage.BeverageId,beverage.Name, beverage.Price,beverage.CreatedAt,beverage.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		//encode the beverage struct into json and write it to the response writer
		json.NewEncoder(w).Encode(beverage)
	}
}

//method to update a beverage
func UpdateBeverage(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract beverage_id from URL path
        vars := mux.Vars(r)
        beverageId := vars["beverage_id"]

        // Fetch the existing record from the database
        var existingBeverage models.Beverage
        query := "SELECT beverage_id, name,price FROM beverages WHERE beverage_id = ?"
        if err := db.QueryRow(query, beverageId).Scan(
            &existingBeverage.BeverageId,
            &existingBeverage.Name,
            &existingBeverage.Price,
        ); err != nil {
            log.Printf("Error fetching existing pizza type: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }


        // Decode the request body into a new beverage struct
        var updatedBeverage models.Beverage
        if err := json.NewDecoder(r.Body).Decode(&updatedBeverage); err != nil {
            log.Printf("Error decoding request body: %v", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Merge the changes
        if updatedBeverage.Name != "" {
            existingBeverage.Name = updatedBeverage.Name
        }
        if updatedBeverage.Price != 0 {
            existingBeverage.Price = updatedBeverage.Price
		}
        existingBeverage.UpdatedAt = time.Now()

        // Build the update query dynamically
        updateQuery := "UPDATE beverages SET name=?, price=?,updated_at=? WHERE beverage_id=?"
        args := []interface{}{
            existingBeverage.Name,
            existingBeverage.Price,
            existingBeverage.UpdatedAt,
            beverageId,
        }

        // Execute the query and check for errors
        if _, err := db.Exec(updateQuery, args...); err != nil {
            log.Printf("Error executing update query: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Set the response header to application/json
        w.Header().Set("Content-Type", "application/json")
        // Encode the updated beverage struct into json and write it to the response writer
        if err := json.NewEncoder(w).Encode(existingBeverage); err != nil {
            log.Printf("Error encoding response: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}
//method to delete a beverage
func DeleteBeverage(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract beverage_id from the request URL
		vars := mux.Vars(r)
		beverageId := vars["beverage_id"]

		// create a query to delete the beverage from the database
		query := "DELETE FROM beverages WHERE beverage_id = ?"
		// execute the query and check for errors
		if _, err := db.Exec(query, beverageId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// set the response header to application/json
		w.Header().Set("Content-Type", "application/json")
		// write a success message to the response writer
		json.NewEncoder(w).Encode(map[string]string{"message": "Beverage deleted successfully"})
	}
}