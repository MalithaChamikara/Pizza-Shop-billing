package models

import "time"
type Invoice struct {
	InvoiceId string `json:"invoice_id"`
	InvoiceDate string `json:"invoice_date"`
	SubTotal float64 `json:"subtotal"`
	Tax float64 `json:"tax"`
	Total float64 `json:"total"`
	CustomerName string `json:"customer_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}