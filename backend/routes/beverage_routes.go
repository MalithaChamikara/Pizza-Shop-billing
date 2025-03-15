package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"piza_shop_billing/backend/controllers"
	"piza_shop_billing/backend/database"
)

func RegisterBeverageRoutes(router *mux.Router) {
	
	router.HandleFunc("/beverages", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            controllers.GetBeverages(database.DB)(w, r)
        case http.MethodPost:
            controllers.CreateBeverage(database.DB)(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "POST")

	//route for updating a topping
	router.HandleFunc("/beverages/{beverage_id}", controllers.UpdateBeverage(database.DB)).Methods("PUT")

	//route for deleting a topping
	router.HandleFunc("/beverages/{beverage_id}", controllers.DeleteBeverage(database.DB)).Methods("DELETE")

}