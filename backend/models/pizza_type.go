package models

import "time"
type PizzaType struct {
	PizzaTypeId string `json:"pizza_type_id"`
	Name		string `json:"name"`
	Size		string `json:"size"`
	BasePrice	float64 `json:"base_price"`
	Description	string `json:"description"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`

}