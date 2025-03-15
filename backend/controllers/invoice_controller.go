package controllers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	"piza_shop_billing/backend/models"
	"time"
)

const DateTimeFormat = "2006-01-02 15:04:05"

//function to return all invoices
func GetInvoices(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        query := "SELECT invoice_id,DATE_FORMAT(invoice_date,'%Y-%m-%d') AS invoice_date, subtotal, tax, total,customer_name FROM invoices"
        results, err := db.Query(query)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer results.Close()

		//create a slice of invoices
        var invoices []models.Invoice
        for results.Next() {
            var invoice models.Invoice
            if err := results.Scan(&invoice.InvoiceId,&invoice.InvoiceDate,&invoice.SubTotal, &invoice.Tax, &invoice.Total,&invoice.CustomerName); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
			//if no error append the invoice to the invoices slice
            invoices = append(invoices, invoice)
        }

		//generate a json response object
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(invoices)
    }
}

//function to create a new invoice
func CreateInvoice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        var invoice models.Invoice
        if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        // Initialize subtotal, tax, and total to 0.00
        invoice.SubTotal = 0.00
        invoice.Tax = 0.00
        invoice.Total = 0.00
		invoice.InvoiceDate = time.Now().Format(DateTimeFormat)
		invoice.UpdatedAt = time.Now()

        query := "INSERT INTO invoices (invoice_date, subtotal, tax, total,customer_name,updated_at) VALUES (?,?,?,?,?,?)"
        _, err := db.Exec(query, invoice.InvoiceDate, invoice.SubTotal, invoice.Tax, invoice.Total,invoice.CustomerName, invoice.UpdatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

    
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
       

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(invoice)
    }
}

//function to update an invoice
func UpdateInvoice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceID := vars["invoice_id"]

		var invoice models.Invoice
		if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		 // Update invoice details (excluding subtotal, tax, and total)
		 query := "UPDATE invoices SET  customer_name=?, updated_at=NOW() WHERE invoice_id=?"
		 _, err := db.Exec(query,invoice.CustomerName, invoiceID)
		 if err != nil {
			 http.Error(w, err.Error(), http.StatusInternalServerError)
			 return
		 }
		 // Calculate subtotal, tax, and total from invoice_items
		 query = `
		 SELECT SUM(quantity * unit_price) AS subtotal
		 FROM invoice_items
		 WHERE invoice_id = ?`

		 var subtotal float64
		 err = db.QueryRow(query, invoiceID).Scan(&subtotal)
		 if err != nil {
			 http.Error(w, err.Error(), http.StatusInternalServerError)
			 return
		 }

		 tax := subtotal * 0.10
		 total := subtotal + tax
		// Update subtotal, tax, and total in invoices table
		query = "UPDATE invoices SET subtotal=?, tax=?, total=? WHERE invoice_id=?"
		_, err = db.Exec(query, subtotal, tax, total, invoiceID)
		if err != nil {
			  http.Error(w, err.Error(), http.StatusInternalServerError)
			  return
		}

		  // Fetch the updated invoice
		  query = "SELECT invoice_id, customer_name, subtotal, tax, total FROM invoices WHERE invoice_id=?"
		  err = db.QueryRow(query, invoiceID).Scan(&invoice.InvoiceId, &invoice.CustomerName, &invoice.SubTotal, &invoice.Tax, &invoice.Total)
		  if err != nil {
			  http.Error(w, err.Error(), http.StatusInternalServerError)
			  return
		  }

		  w.Header().Set("Content-Type", "application/json")
		  json.NewEncoder(w).Encode(invoice)
  	
	}
}

//delete an invoice
func DeleteInvoice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceID := vars["invoice_id"]

		query := "DELETE FROM invoices WHERE invoice_id = ?"
		_, err := db.Exec(query, invoiceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"message": "Invoice deleted successfully"})
	}
}

//function to return all invoice items specific to an invoice
func GetInvoiceItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceID := vars["invoice_id"]

		query := "SELECT invoice_item_id, invoice_id, item_id, quantity, unit_price FROM invoice_items WHERE invoice_id = ?"
		results, err := db.Query(query, invoiceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer results.Close()

		var invoiceItems []models.InvoiceItem
		for results.Next() {
			var invoiceItem models.InvoiceItem
			if err := results.Scan(&invoiceItem.InvoiceItemId, &invoiceItem.InvoiceId, &invoiceItem.ItemId, &invoiceItem.Quantity, &invoiceItem.UnitPrice); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			invoiceItems = append(invoiceItems, invoiceItem)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoiceItems)
	}
}

//function to create a new invoice item
func CreateInvoiceItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceID := vars["invoice_id"]

		var invoiceItem models.InvoiceItem
		if err := json.NewDecoder(r.Body).Decode(&invoiceItem); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO invoice_items (invoice_id, item_id, quantity, unit_price) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(query, invoiceID, invoiceItem.ItemId, invoiceItem.Quantity, invoiceItem.UnitPrice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(invoiceItem)
	}
}

//funtion to update an invoice item
func UpdateInvoiceItem(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        itemID := vars["invoice_item_id"]

        var item models.InvoiceItem
        if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        query := "UPDATE invoice_items SET item_id=?, quantity=?, unit_price=?,  WHERE invoice_item_id=?"
        _, err := db.Exec(query, item.ItemId, item.Quantity, item.UnitPrice,itemID)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(item)
    }
}

//function to delete an invoice item
func DeleteInvoiceItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceItemID := vars["invoice_item_id"]

		query := "DELETE FROM invoice_items WHERE   invoice_item_id = ?"
		_, err := db.Exec(query,invoiceItemID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Invoice item deleted successfully"})
	}
}

//function to generate a printable invoice
func GeneratePrintableInvoice(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		invoiceID := vars["invoice_id"]

		//fetch the invoice details
		
		query := "SELECT invoice_id,subtotal, tax, total FROM invoices WHERE invoice_id = ?"
		results, err := db.Query(query, invoiceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer results.Close()

		var invoice models.Invoice
		if results.Next() {
			if err := results.Scan(&invoice.InvoiceId, &invoice.SubTotal, &invoice.Tax, &invoice.Total); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		//fetch the invoice items

		query = "SELECT item_id, quantity, unit_price FROM invoice_items WHERE invoice_id = ?"
		results, err = db.Query(query, invoiceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer results.Close()

		var invoiceItems []models.InvoiceItem
		for results.Next() {
			var invoiceItem models.InvoiceItem
			if err := results.Scan(&invoiceItem.ItemId, &invoiceItem.Quantity, &invoiceItem.UnitPrice); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			invoiceItems = append(invoiceItems, invoiceItem)
		}

		//create response object
		response := struct {
			Invoice models.Invoice `json:"invoice"`
			InvoiceItems []models.InvoiceItem `json:"invoice_items"`
		}{
			Invoice: invoice,
			InvoiceItems: invoiceItems,
		}

		

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}