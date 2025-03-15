package main

import (
    
    "net/http"
    "log"
    "piza_shop_billing/backend/routes"
    "piza_shop_billing/backend/database"
    "github.com/rs/cors"
    
)

func main() {

    
    database.Connect()
    defer database.DB.Close()
    log.Println("Application connected to the database")

   // Register the pizza routes
   router := routes.RegisterPizzaRoutes()

   // Register the topping routes
   routes.RegisterToppingRoutes(router)

    // Register the beverage routes
    routes.RegisterBeverageRoutes(router)

    // Register the invoice routes
    routes.RegisterInvoiceRoutes(router)

   // Define the root path
   router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       w.Header().Set("Content-Type", "application/json")
       w.Write([]byte(`{"message": "Welcome to the Pizza Shop Billing API"}`))
   })

    // Enable CORS with default settings (allowing all origins)
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Allows all origins
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"}, 
    })

    // start the server on port 8080
    if err := http.ListenAndServe(":8080", c.Handler(router)); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }

   
}


