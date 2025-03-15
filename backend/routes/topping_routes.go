package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"piza_shop_billing/backend/controllers"
	"piza_shop_billing/backend/database"
)

func RegisterToppingRoutes(router *mux.Router) {
	//create a new router with the mux package
	router.HandleFunc("/toppings", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            controllers.GetToppings(database.DB)(w, r)
        case http.MethodPost:
            controllers.CreateTopping(database.DB)(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "POST")

	//route for updating a topping
	router.HandleFunc("/toppings/{topping_id}", controllers.UpdateTopping(database.DB)).Methods("PUT")

	//route for deleting a topping
	router.HandleFunc("/toppings/{topping_id}", controllers.DeleteTopping(database.DB)).Methods("DELETE")

}