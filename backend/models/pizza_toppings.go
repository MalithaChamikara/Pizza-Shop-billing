package models

type PizzaTopping struct {
	PizzaToppingId int `json:"pizza_topping_id"`
	PizzaTypeId string `json:"pizza_type_id"`
	ToppingId string `json:"topping_id"`
}
