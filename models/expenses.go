package models

import "time"

type GetExpense struct {
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Comment  string  `json:"comment"`
}

type Expense struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
