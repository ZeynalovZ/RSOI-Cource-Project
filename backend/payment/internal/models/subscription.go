package models

type Subscription struct {
	Id 			int `json:"id" db:"id"`
	Name 		string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Price 		int `json:"price" db:"price"`
	Status 		int `json:"status" db:"status"`
}

const (
	User = iota
	PrimeUser
)