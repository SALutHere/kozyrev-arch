package model

// Part - деталь со склада
type Part struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Quantity int     `json:"quantity"`
	Weight   float64 `json:"weight"`
}
