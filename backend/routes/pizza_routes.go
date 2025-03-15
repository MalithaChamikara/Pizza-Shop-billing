package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"piza_shop_billing/backend/controllers"
	"piza_shop_billing/backend/database"

)
func RegisterPizzaRoutes()  *mux.Router {


	//create a new router with the mux package
	router := mux.NewRouter()
	router.HandleFunc("/pizzas", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            controllers.GetPizzaTypes(database.DB)(w, r)
        case http.MethodPost:
            controllers.CreatePizzaType(database.DB)(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "POST")

	//route for updating a pizza type
	router.HandleFunc("/pizzas/{pizza_type_id}", controllers.UpdatePizzaType(database.DB)).Methods("PUT")

	//route for deleting a pizza type
	router.HandleFunc("/pizzas/{pizza_type_id}", controllers.DeletePizzaType(database.DB)).Methods("DELETE")

	//route for linking a pizza type and topping
	router.HandleFunc("/pizzas/{pizza_type_id}/toppings", controllers.LinkPizzaTopping(database.DB)).Methods("POST")

	//route for get topping name for specific pizza type
	router.HandleFunc("/pizzas/{pizza_type_id}/toppings", controllers.GetToppingsByPizzaType(database.DB)).Methods("GET")

	return router
	
}