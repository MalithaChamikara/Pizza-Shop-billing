package routes

import (
    "github.com/gorilla/mux"
    "piza_shop_billing/backend/controllers"
    "piza_shop_billing/backend/database"
)

func RegisterInvoiceRoutes(router *mux.Router) {
    router.HandleFunc("/invoices", controllers.GetInvoices(database.DB)).Methods("GET")
    router.HandleFunc("/invoices", controllers.CreateInvoice(database.DB)).Methods("POST")
    router.HandleFunc("/invoices/{invoice_id}", controllers.UpdateInvoice(database.DB)).Methods("PUT")
    router.HandleFunc("/invoices/{invoice_id}", controllers.DeleteInvoice(database.DB)).Methods("DELETE")

    router.HandleFunc("/invoices/{invoice_id}/items", controllers.GetInvoiceItems(database.DB)).Methods("GET")
    router.HandleFunc("/invoices/{invoice_id}/items", controllers.CreateInvoiceItem(database.DB)).Methods("POST")
    router.HandleFunc("/invoices/items/{invoice_item_id}", controllers.UpdateInvoiceItem(database.DB)).Methods("PUT")
    router.HandleFunc("/invoices/items/{invoice_item_id}", controllers.DeleteInvoiceItem(database.DB)).Methods("DELETE")

    router.HandleFunc("/invoices/{invoice_id}/print", controllers.GeneratePrintableInvoice(database.DB)).Methods("GET")
}