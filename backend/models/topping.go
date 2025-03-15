package models

import "time"
type Topping struct {
	ToppingId string `json:"topping_id"`
	Name		string `json:"name"`
	Price		float64 `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}