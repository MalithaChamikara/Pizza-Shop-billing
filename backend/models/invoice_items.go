package models

type InvoiceItem struct {
    InvoiceItemId int     `json:"invoice_item_id"`
    InvoiceId     int     `json:"invoice_id"`
    ItemId        string  `json:"item_id,omitempty"`
    Quantity      int     `json:"quantity"`
    UnitPrice     float64 `json:"unit_price"`
    
}